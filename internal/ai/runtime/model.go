package runtime

import (
	_ "embed"
	"encoding/binary"
	"errors"
	"unsafe"
)

// ---- embedded model ----

//go:embed model.bin
var ModelData []byte

// ---- header ----

type ModelHeader struct {
	Magic                 uint32
	Version               uint32
	VocabSize             uint32
	HiddenSize            uint32
	IntermediateSize      uint32
	NumLayers             uint32
	MaxPositionEmbeddings uint32
	TokenTypeVocabSize    uint32
}

const (
	modelMagic   = 0x4C414253 // "LABS"
	modelVersion = 1
)

// ---- tensor descriptors ----

// Linear int8 + scales + bias
type LinearInt8 struct {
	InFeatures  int
	OutFeatures int

	Weights []int8    // len = out * in
	Scales  []float32 // len = out
	Bias    []float32 // len = out
}

// LayerNorm float32
type LayerNorm struct {
	Gamma []float32 // len = hidden
	Beta  []float32 // len = hidden
}

// Embeddings block
type Embeddings struct {
	Word      []float32 // [vocab][hidden]
	Position  []float32 // [max_pos][hidden]
	TokenType []float32 // [token_types][hidden]
	Gamma     []float32 // [hidden]
	Beta      []float32 // [hidden]
}

// One encoder layer weights
type EncoderLayer struct {
	Q        LinearInt8
	K        LinearInt8
	V        LinearInt8
	O        LinearInt8
	FFNUp    LinearInt8
	FFNDown  LinearInt8
	AttnLN   LayerNorm
	OutputLN LayerNorm
}

// ---- main model struct ----

type Model struct {
	H ModelHeader

	Embeddings Embeddings
	Layers     []EncoderLayer

	// cached sizes
	hidden   int
	inter    int
	numHeads int
	headDim  int
}

// ---- public API ----

func LoadModel(data []byte) (*Model, error) {
	if len(data) < 32 {
		return nil, errors.New("model: file too small")
	}

	var h ModelHeader
	h.Magic = binary.LittleEndian.Uint32(data[0:4])
	h.Version = binary.LittleEndian.Uint32(data[4:8])
	h.VocabSize = binary.LittleEndian.Uint32(data[8:12])
	h.HiddenSize = binary.LittleEndian.Uint32(data[12:16])
	h.IntermediateSize = binary.LittleEndian.Uint32(data[16:20])
	h.NumLayers = binary.LittleEndian.Uint32(data[20:24])
	h.MaxPositionEmbeddings = binary.LittleEndian.Uint32(data[24:28])
	h.TokenTypeVocabSize = binary.LittleEndian.Uint32(data[28:32])

	if h.Magic != modelMagic {
		return nil, errors.New("model: invalid magic")
	}
	if h.Version != modelVersion {
		return nil, errors.New("model: unsupported version")
	}

	// bge-small-en-v1.5: hidden=384, num_heads=12, head_dim=32
	m := &Model{
		H:        h,
		hidden:   int(h.HiddenSize),
		inter:    int(h.IntermediateSize),
		numHeads: 12,
		headDim:  int(h.HiddenSize) / 12,
	}

	offset := 32
	var err error

	offset, err = m.parseEmbeddings(data, offset)
	if err != nil {
		return nil, err
	}

	offset, err = m.parseLayers(data, offset)
	if err != nil {
		return nil, err
	}

	// optional: sanity check end offset == len(data)
	if offset != len(data) {
		// not fatal, but useful for debugging
	}

	return m, nil
}

// ---- parsing helpers ----

func (m *Model) parseEmbeddings(data []byte, offset int) (int, error) {
	vocab := int(m.H.VocabSize)
	hidden := m.hidden
	maxPos := int(m.H.MaxPositionEmbeddings)
	tokenTypes := int(m.H.TokenTypeVocabSize)

	needFloats := vocab*hidden + maxPos*hidden + tokenTypes*hidden + 2*hidden
	needBytes := needFloats * 4
	if offset+needBytes > len(data) {
		return 0, errors.New("model: embeddings truncated")
	}

	buf := data[offset : offset+needBytes]
	f := bytesToFloat32(buf)

	pos := 0

	m.Embeddings.Word = f[pos : pos+vocab*hidden]
	pos += vocab * hidden

	m.Embeddings.Position = f[pos : pos+maxPos*hidden]
	pos += maxPos * hidden

	m.Embeddings.TokenType = f[pos : pos+tokenTypes*hidden]
	pos += tokenTypes * hidden

	m.Embeddings.Gamma = f[pos : pos+hidden]
	pos += hidden

	m.Embeddings.Beta = f[pos : pos+hidden]
	pos += hidden

	return offset + needBytes, nil
}

func (m *Model) parseLayers(data []byte, offset int) (int, error) {
	layers := int(m.H.NumLayers)
	m.Layers = make([]EncoderLayer, layers)

	for i := 0; i < layers; i++ {
		var err error
		// Q
		offset, err = parseLinearInt8(data, offset, &m.Layers[i].Q)
		if err != nil {
			return 0, err
		}
		// K
		offset, err = parseLinearInt8(data, offset, &m.Layers[i].K)
		if err != nil {
			return 0, err
		}
		// V
		offset, err = parseLinearInt8(data, offset, &m.Layers[i].V)
		if err != nil {
			return 0, err
		}
		// O
		offset, err = parseLinearInt8(data, offset, &m.Layers[i].O)
		if err != nil {
			return 0, err
		}
		// FFN_up
		offset, err = parseLinearInt8(data, offset, &m.Layers[i].FFNUp)
		if err != nil {
			return 0, err
		}
		// FFN_down
		offset, err = parseLinearInt8(data, offset, &m.Layers[i].FFNDown)
		if err != nil {
			return 0, err
		}
		// Attn LN
		offset, err = parseLayerNorm(data, offset, &m.Layers[i].AttnLN, m.hidden)
		if err != nil {
			return 0, err
		}
		// Output LN
		offset, err = parseLayerNorm(data, offset, &m.Layers[i].OutputLN, m.hidden)
		if err != nil {
			return 0, err
		}
	}

	return offset, nil
}

func parseLinearInt8(data []byte, offset int, dst *LinearInt8) (int, error) {
	if offset+12 > len(data) {
		return 0, errors.New("model: linear header truncated")
	}
	typ := binary.LittleEndian.Uint32(data[offset : offset+4])
	if typ != 3 {
		return 0, errors.New("model: expected linear block type=3")
	}
	inF := int(binary.LittleEndian.Uint32(data[offset+4 : offset+8]))
	outF := int(binary.LittleEndian.Uint32(data[offset+8 : offset+12]))
	offset += 12

	weightsBytes := inF * outF
	scalesBytes := outF * 4
	biasBytes := outF * 4

	if offset+weightsBytes+scalesBytes+biasBytes > len(data) {
		return 0, errors.New("model: linear data truncated")
	}

	dst.InFeatures = inF
	dst.OutFeatures = outF

	dst.Weights = bytesToInt8(data[offset : offset+weightsBytes])
	offset += weightsBytes

	dst.Scales = bytesToFloat32(data[offset : offset+scalesBytes])
	offset += scalesBytes

	dst.Bias = bytesToFloat32(data[offset : offset+biasBytes])
	offset += biasBytes

	return offset, nil
}

func parseLayerNorm(data []byte, offset int, dst *LayerNorm, hidden int) (int, error) {
	if offset+4 > len(data) {
		return 0, errors.New("model: layernorm header truncated")
	}
	typ := binary.LittleEndian.Uint32(data[offset : offset+4])
	if typ != 2 {
		return 0, errors.New("model: expected layernorm block type=2")
	}
	offset += 4

	needFloats := hidden * 2
	needBytes := needFloats * 4
	if offset+needBytes > len(data) {
		return 0, errors.New("model: layernorm data truncated")
	}

	buf := bytesToFloat32(data[offset : offset+needBytes])
	dst.Gamma = buf[:hidden]
	dst.Beta = buf[hidden : 2*hidden]

	return offset + needBytes, nil
}

// ---- byte views ----

func bytesToFloat32(b []byte) []float32 {
	hdr := (*[3]uintptr)(unsafe.Pointer(&b))
	var s []float32
	sh := (*[3]uintptr)(unsafe.Pointer(&s))
	sh[0] = hdr[0]
	sh[1] = uintptr(len(b) / 4)
	sh[2] = uintptr(len(b) / 4)
	return s
}

func bytesToInt8(b []byte) []int8 {
	hdr := (*[3]uintptr)(unsafe.Pointer(&b))
	var s []int8
	sh := (*[3]uintptr)(unsafe.Pointer(&s))
	sh[0] = hdr[0]
	sh[1] = uintptr(len(b))
	sh[2] = uintptr(len(b))
	return s
}
