package epub2

import (
	"archive/zip"
	"embed"
	"io"
	"text/template"
	"time"

	"github.com/google/uuid"
)

type EPUB struct {
	UUID     string
	Lang     string
	Title    string
	Metadata string
	Manifest string
	Spine    string
	Toc      []TOC

	zw   *zip.Writer
	tmpl *template.Template
}

type TOC struct {
	Id    string
	Order int
	Text  string
	Src   string
	Depth int
}

//go:embed assets/*
var assets embed.FS

func New(wc io.WriteCloser) (*EPUB, error) {
	epub := &EPUB{
		UUID:     uuid.New().String(),
		Manifest: `<item id="css" href="main.css" media-type="text/css" />` + "\n",
		Toc:      make([]TOC, 0),
	}

	epub.zw = zip.NewWriter(wc)
	for _, f := range []string{
		"mimetype",
		"META-INF/container.xml",
		"OEBPS/main.css",
	} {
		var header *zip.FileHeader
		if f == "mimetype" {
			header = &zip.FileHeader{
				Name:   f,
				Method: zip.Store,
			}
		} else {
			header = &zip.FileHeader{
				Name:     f,
				Method:   zip.Deflate,
				Modified: time.Now(),
			}
		}
		dst, err := epub.zw.CreateHeader(header)
		if err != nil {
			panic(err)
		}
		content, err := assets.ReadFile("assets/files/" + f)
		if err != nil {
			panic(err)
		}
		_, err = dst.Write(content)
		if err != nil {
			panic(err)
		}
	}
	epub.tmpl = template.Must(template.New("").ParseFS(assets, "assets/tmpl/*.tmpl"))

	return epub, nil
}

func (e *EPUB) execTemplate(file, name string, data any) error {
	header := &zip.FileHeader{
		Name:     file,
		Method:   zip.Deflate,
		Modified: time.Now(),
	}
	w, err := e.zw.CreateHeader(header)
	if err != nil {
		return err
	}
	return e.tmpl.ExecuteTemplate(w, name, data)
}

func (e *EPUB) Close() error {
	return e.zw.Close()
}
