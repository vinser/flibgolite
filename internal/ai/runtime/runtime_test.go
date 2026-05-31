package runtime_test

import (
	"encoding/json"
	"math"
	"os"
	"testing"

	"github.com/vinser/flibgolite/internal/ai/runtime"
)

func TestEmbedTokens(t *testing.T) {
	// 1) Load model
	m, err := runtime.LoadModel(runtime.ModelData)
	if err != nil {
		t.Fatalf("LoadModel failed: %v", err)
	}

	// 2) Fixed tokens (CLS, "this", SEP)
	ids := []int32{101, 2023, 102}

	// 3) Run embedding
	emb := m.EmbedTokens(ids)

	if len(emb) != int(m.H.HiddenSize) {
		t.Fatalf("expected embedding size %d, got %d", m.H.HiddenSize, len(emb))
	}

	// 4) Check L2 normalization
	norm := float32(0)
	for _, v := range emb {
		norm += v * v
	}
	norm = float32(math.Sqrt(float64(norm)))

	if math.Abs(float64(norm-1.0)) > 1e-4 {
		t.Fatalf("embedding not normalized, norm=%f", norm)
	}

	t.Logf("EmbedTokens OK, norm=%f", norm)
}

func TestCompareWithPyTorch(t *testing.T) {
	m, err := runtime.LoadModel(runtime.ModelData)
	if err != nil {
		t.Fatal(err)
	}

	tokens := []int32{101, 7592, 2088, 102} // [CLS] hello world [SEP]
	v := m.EmbedTokens(tokens)

	// ---- Диагностика ref.txt ----
	b, err := os.ReadFile("/home/vinser/Dev/flibgolite/tools/ai/model/convert/ref.txt")
	if err != nil {
		t.Fatalf("cannot read ref.txt: %v", err)
	}
	t.Logf("ref.txt size = %d bytes", len(b))

	var ref []float32
	if err := json.Unmarshal(b, &ref); err != nil {
		t.Fatalf("json error: %v", err)
	}
	t.Logf("ref length = %d", len(ref))
	// --------------------------------

	cos := cosine(v, ref)
	t.Logf("cosine = %.6f", cos)

	if cos < 0.999 {
		t.Fatalf("cosine too low: %.6f", cos)
	}
}

func loadRef(path string) []float32 {
	b, _ := os.ReadFile(path)
	var arr []float32
	json.Unmarshal(b, &arr)
	return arr
}

func cosine(a, b []float32) float32 {
	var dot, na, nb float32
	for i := range a {
		dot += a[i] * b[i]
		na += a[i] * a[i]
		nb += b[i] * b[i]
	}
	return dot / float32(math.Sqrt(float64(na*nb)))
}
