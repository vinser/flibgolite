package runtime

import "math"

// ---- Linear INT8 ----

func (l *LinearInt8) Forward(dst, src []float32) {
	// dst: outFeatures
	// src: inFeatures

	out := l.OutFeatures
	in := l.InFeatures

	for o := 0; o < out; o++ {
		sum := float32(0)

		wRow := l.Weights[o*in : o*in+in]
		scale := l.Scales[o]
		bias := l.Bias[o]

		for i := 0; i < in; i++ {
			sum += float32(wRow[i]) * src[i]
		}

		dst[o] = sum*scale + bias
	}
}

// ---- LayerNorm ----

func (ln *LayerNorm) Forward(dst, src []float32) {
	n := len(src)

	mean := float32(0)
	for i := 0; i < n; i++ {
		mean += src[i]
	}
	mean /= float32(n)

	variance := float32(0)
	for i := 0; i < n; i++ {
		d := src[i] - mean
		variance += d * d
	}
	variance /= float32(n)

	const eps = 1e-5 // match PyTorch LayerNorm eps
	invStd := 1.0 / float32(math.Sqrt(float64(variance)+eps))

	for i := 0; i < n; i++ {
		norm := (src[i] - mean) * invStd
		dst[i] = norm*ln.Gamma[i] + ln.Beta[i]
	}
}

// ---- GELU ----

func gelu(x float32) float32 {
	const c = 0.044715
	return 0.5 * x * (1 + float32(math.Tanh(float64(0.79788456*(x+c*x*x*x)))))
}

// ---- Self-Attention ----

type AttentionScratch struct {
	// [seq][hidden]
	Q       [][]float32
	K       [][]float32
	V       [][]float32
	Context [][]float32
	AttnOut [][]float32

	// [heads][seq][headDim]
	QHeads       [][][]float32
	KHeads       [][][]float32
	VHeads       [][][]float32
	ContextHeads [][][]float32

	// [seq][seq]
	Scores [][]float32
}

// SelfAttention runs self-attention for a single encoder layer.
func (l *EncoderLayer) SelfAttention(
	out [][]float32,
	x [][]float32,
	tmp *AttentionScratch,
	numHeads int,
	headDim int,
) {
	seqLen := len(x)
	hidden := len(x[0])

	// 1) Q, K, V projections
	for t := 0; t < seqLen; t++ {
		l.Q.Forward(tmp.Q[t], x[t])
		l.K.Forward(tmp.K[t], x[t])
		l.V.Forward(tmp.V[t], x[t])
	}

	// 2) split into heads
	splitHeads(tmp.QHeads, tmp.Q, numHeads, headDim)
	splitHeads(tmp.KHeads, tmp.K, numHeads, headDim)
	splitHeads(tmp.VHeads, tmp.V, numHeads, headDim)

	// 3) attention per head
	scale := float32(1.0 / math.Sqrt(float64(headDim)))
	for h := 0; h < numHeads; h++ {
		attnHead(
			tmp.ContextHeads[h],
			tmp.QHeads[h],
			tmp.KHeads[h],
			tmp.VHeads[h],
			tmp.Scores,
			scale,
		)
	}

	// 4) merge heads
	mergeHeads(tmp.Context, tmp.ContextHeads, numHeads, headDim)

	// 5) output projection + residual + LayerNorm
	for t := 0; t < seqLen; t++ {
		l.O.Forward(tmp.AttnOut[t], tmp.Context[t])
		for i := 0; i < hidden; i++ {
			tmp.AttnOut[t][i] += x[t][i]
		}
		l.AttnLN.Forward(out[t], tmp.AttnOut[t])
	}
}

// splitHeads: [seq][hidden] -> [heads][seq][headDim]
func splitHeads(dst [][][]float32, src [][]float32, numHeads, headDim int) {
	seqLen := len(src)
	for h := 0; h < numHeads; h++ {
		for t := 0; t < seqLen; t++ {
			base := h * headDim
			copy(dst[h][t], src[t][base:base+headDim])
		}
	}
}

// mergeHeads: [heads][seq][headDim] -> [seq][hidden]
func mergeHeads(dst [][]float32, src [][][]float32, numHeads, headDim int) {
	seqLen := len(dst)
	for t := 0; t < seqLen; t++ {
		pos := 0
		for h := 0; h < numHeads; h++ {
			copy(dst[t][pos:pos+headDim], src[h][t])
			pos += headDim
		}
	}
}

// attnHead computes attention for a single head.
func attnHead(
	out [][]float32,
	q [][]float32,
	k [][]float32,
	v [][]float32,
	scores [][]float32,
	scale float32,
) {
	seqLen := len(q)
	headDim := len(q[0])

	// scores = Q * K^T * scale
	for i := 0; i < seqLen; i++ {
		for j := 0; j < seqLen; j++ {
			sum := float32(0)
			qi := q[i]
			kj := k[j]
			for d := 0; d < headDim; d++ {
				sum += qi[d] * kj[d]
			}
			scores[i][j] = sum * scale
		}
	}

	// softmax over j
	for i := 0; i < seqLen; i++ {
		row := scores[i]

		maxVal := row[0]
		for j := 1; j < seqLen; j++ {
			if row[j] > maxVal {
				maxVal = row[j]
			}
		}

		sum := float32(0)
		for j := 0; j < seqLen; j++ {
			e := float32(math.Exp(float64(row[j] - maxVal)))
			row[j] = e
			sum += e
		}

		inv := 1.0 / sum
		for j := 0; j < seqLen; j++ {
			row[j] *= inv
		}
	}

	// out = scores * V
	for i := 0; i < seqLen; i++ {
		oi := out[i]
		for d := 0; d < headDim; d++ {
			oi[d] = 0
		}
		for j := 0; j < seqLen; j++ {
			w := scores[i][j]
			vj := v[j]
			for d := 0; d < headDim; d++ {
				oi[d] += w * vj[d]
			}
		}
	}
}
