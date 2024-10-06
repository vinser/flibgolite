package fb2

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func BenchmarkParseFB2(b *testing.B) {
	a, err := os.ReadFile("testdata/303934.fb2")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	b.Run("encoding/xml", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fb2, err := ParseFB2(io.NopCloser(bytes.NewReader(a)))
			if err != nil {
				b.Fatal(err)
			}
			_ = fb2
		}
	})
	b.Run("gosax", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fb2, err := ParseFB2Gosax(io.NopCloser(bytes.NewReader(a)))
			if err != nil {
				b.Fatal(err)
			}
			_ = fb2
		}
	})
}
