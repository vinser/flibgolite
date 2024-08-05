package epub2

import (
	"archive/zip"
	"embed"
	"io"
	"io/fs"
	"text/template"

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
	files, err := fs.Sub(assets, "assets/files")
	if err != nil {
		panic(err)
	}
	epub.zw.AddFS(files)
	epub.tmpl = template.Must(template.New("").ParseFS(assets, "assets/tmpl/*.tmpl"))

	return epub, nil
}

func (e *EPUB) execTemplate(file, name string, data any) error {
	w, err := e.zw.Create(file)
	if err != nil {
		return err
	}
	return e.tmpl.ExecuteTemplate(w, name, data)
}

func (e *EPUB) Close() error {
	return e.zw.Close()
}
