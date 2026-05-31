package runtime

import (
	"math"
)

// EmbedTokens runs the full encoder on a sequence of token IDs
// and returns a normalized embedding of size hidden.
func (m *Model) EmbedTokens(ids []int32) []float32 {
	seqLen := len(ids)
	hidden := m.hidden

	// 1) allocate sequence buffers
	x := make2D(seqLen, hidden)
	xNext := make2D(seqLen, hidden)

	// 2) embeddings + LayerNorm
	m.applyEmbeddings(x, ids)

	// 3) attention scratch
	attnScratch := NewAttentionScratch(seqLen, hidden, m.numHeads, m.headDim)

	// 4) encoder layers
	ffnHidden := make([]float32, m.inter)

	for li := 0; li < int(m.H.NumLayers); li++ {
		layer := &m.Layers[li]

		// self-attention
		layer.SelfAttention(xNext, x, attnScratch, m.numHeads, m.headDim)

		// swap x/xNext
		x, xNext = xNext, x

		// FFN + residual + LayerNorm (per token)
		for t := 0; t < seqLen; t++ {
			// FFN_up: hidden -> intermediate
			layer.FFNUp.Forward(ffnHidden, x[t])

			for i := 0; i < m.inter; i++ {
				ffnHidden[i] = gelu(ffnHidden[i])
			}

			// FFN_down: intermediate -> hidden
			layer.FFNDown.Forward(xNext[t], ffnHidden)

			// residual
			for i := 0; i < hidden; i++ {
				xNext[t][i] += x[t][i]
			}

			// LayerNorm
			layer.OutputLN.Forward(xNext[t], xNext[t])
		}

		// swap again
		x, xNext = xNext, x
	}

	// 5) CLS pooling (take token 0)
	emb := make([]float32, hidden)
	copy(emb, x[0])

	// 6) L2 normalization
	norm := float32(0)
	for i := 0; i < hidden; i++ {
		norm += emb[i] * emb[i]
	}
	norm = float32(math.Sqrt(float64(norm)) + 1e-12)
	invNorm := 1.0 / norm
	for i := 0; i < hidden; i++ {
		emb[i] *= invNorm
	}

	return emb
}

// applyEmbeddings fills x[seq][hidden] from token ids and applies embeddings LayerNorm.
func (m *Model) applyEmbeddings(x [][]float32, ids []int32) {
	seqLen := len(ids)
	hidden := m.hidden

	vocab := int(m.H.VocabSize)
	maxPos := int(m.H.MaxPositionEmbeddings)

	// temporary buffer for pre-LN sum
	tmp := make([]float32, hidden)

	ln := LayerNorm{
		Gamma: m.Embeddings.Gamma,
		Beta:  m.Embeddings.Beta,
	}

	for t := 0; t < seqLen; t++ {
		id := int(ids[t])
		if id < 0 || id >= vocab {
			id = 0 // [UNK] fallback
		}
		pos := t
		if pos >= maxPos {
			pos = maxPos - 1
		}

		// sum word + position + token_type(0)
		for i := 0; i < hidden; i++ {
			w := m.Embeddings.Word[id*hidden+i]
			p := m.Embeddings.Position[pos*hidden+i]
			tt := m.Embeddings.TokenType[0*hidden+i]
			tmp[i] = w + p + tt
		}

		ln.Forward(x[t], tmp)
	}
}

// ---- scratch helpers ----

func NewAttentionScratch(seqLen, hidden, numHeads, headDim int) *AttentionScratch {
	s := &AttentionScratch{
		Q:            make2D(seqLen, hidden),
		K:            make2D(seqLen, hidden),
		V:            make2D(seqLen, hidden),
		Context:      make2D(seqLen, hidden),
		AttnOut:      make2D(seqLen, hidden),
		QHeads:       make3D(numHeads, seqLen, headDim),
		KHeads:       make3D(numHeads, seqLen, headDim),
		VHeads:       make3D(numHeads, seqLen, headDim),
		ContextHeads: make3D(numHeads, seqLen, headDim),
		Scores:       make2D(seqLen, seqLen),
	}
	return s
}

func make2D(n, m int) [][]float32 {
	x := make([][]float32, n)
	for i := range x {
		x[i] = make([]float32, m)
	}
	return x
}

func make3D(a, b, c int) [][][]float32 {
	x := make([][][]float32, a)
	for i := range x {
		x[i] = make2D(b, c)
	}
	return x
}
