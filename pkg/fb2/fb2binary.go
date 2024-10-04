package fb2

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"io"
	"os"
	"path"
	"strings"

	"github.com/orisano/gosax"
	"github.com/vinser/flibgolite/pkg/model"
	"github.com/vinser/u8xml"
)

// Any binary data that is required for the presentation of this book in base64 format.
// Currently only images are used.
type Binary struct {
	ID          string
	ContentType string
	Content     []byte
}

func GetCoverPageBinary(coverLink string, rc io.ReadCloser) (*Binary, error) {
	u8r, err := u8xml.NewReader(rc)
	if err != nil {
		return nil, err
	}
	r := gosax.NewReader(u8r)
	b := &Binary{
		ID: strings.TrimPrefix(coverLink, "#"),
	}
	for {
		e, err := r.Event()
		if err != nil {
			return nil, err
		}
		if e.Type() == gosax.EventEOF {
			return nil, ErrUnexpectedEOF
		}

		name, _ := gosax.Name(e.Bytes)
		switch e.Type() {
		case gosax.EventStart:
			if string(name) == "binary" {
				if getAttr(e.Bytes, "id") == b.ID {
					b.ContentType = getAttr(e.Bytes, "content-type")
					b.Content = []byte(getText(r))
					return b, nil
				}
			}
		}
	}
}

func (b *Binary) String() string {
	return fmt.Sprintf(
		`CoverPage ----
  ID: %s
  Content-type: %s
================================
%#v
===========(100)================
`, b.ID, b.ContentType, b.Content[:99])
}

func GetCoverImage(stock string, book *model.Book) (image.Image, error) {
	var rc io.ReadCloser
	if book.Archive == "" {
		rc, _ = os.Open(path.Join(stock, book.File))
	} else {
		zr, err := zip.OpenReader(path.Join(stock, book.Archive))
		if err != nil {
			return nil, err
		}
		defer zr.Close()
		for _, file := range zr.File {
			if file.Name == book.File {
				rc, _ = file.Open()
				break
			}
		}
	}
	defer rc.Close()
	b, err := GetCoverPageBinary(book.Cover, rc)
	if err != nil {
		return nil, err
	}
	br := base64.NewDecoder(base64.StdEncoding, bytes.NewReader(b.Content))
	img, _, err := image.Decode(br)
	if err != nil {
		return nil, err
	}
	return img, nil
}
