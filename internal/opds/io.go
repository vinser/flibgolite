package opds

import (
	"bytes"
	"io"
	"net/http"
)

type ResponseWriteCloser struct {
	http.ResponseWriter
}

func NewWriteCloser(w http.ResponseWriter) *ResponseWriteCloser {
	return &ResponseWriteCloser{
		ResponseWriter: w,
	}
}

func (w ResponseWriteCloser) Write(b []byte) (int, error) {
	return w.ResponseWriter.Write(b)
}

func (w ResponseWriteCloser) Close() error {
	return nil
}

type BufferedReadSeekCloser struct {
	io.ReadSeeker
}

func NewReadSeekCloser(r io.ReadCloser) (*BufferedReadSeekCloser, error) {
	if rs, ok := r.(io.ReadSeeker); ok {
		return &BufferedReadSeekCloser{
			ReadSeeker: rs,
		}, nil
	}
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	rs := bytes.NewReader(b)

	return &BufferedReadSeekCloser{
		ReadSeeker: rs,
	}, nil
}

func (r BufferedReadSeekCloser) Close() error {
	return nil
}
