package fb2

import (
	"bytes"
	"fmt"
	"io"

	"github.com/orisano/gosax"
	"github.com/vinser/flibgolite/pkg/conv/epub2"
	"github.com/vinser/flibgolite/pkg/database"
	"github.com/vinser/flibgolite/pkg/genres"
	"github.com/vinser/flibgolite/pkg/rlog"
	"github.com/vinser/u8xml"
)

type FB2Parser struct {
	BookId int64
	DB     *database.DB
	GT     *genres.GenresTree
	LOG    *rlog.Log
	RC     io.ReadSeekCloser
	*gosax.Reader

	chapterNum int
	parent     *tagStack
}

func (p *FB2Parser) Restart() error {
	if _, err := p.RC.Seek(0, io.SeekStart); err != nil {
		return err
	}
	u8r, err := u8xml.NewReader(p.RC)
	if err != nil {
		return err
	}
	p.Reader = gosax.NewReader(u8r)
	return nil
}

func (p *FB2Parser) links() (map[string]string, error) {
	p.chapterNum = 0

	var (
		links        = map[string]string{}
		bodyName     string
		bodyNum      int
		itemName     string
		sectionDepth int

		updatePage = func() {
			sectionDepth = 0
			if bodyName == "chapter" {
				p.chapterNum++
				itemName = fmt.Sprintf("%s_%d", bodyName, p.chapterNum)
			} else {
				itemName = bodyName
			}
		}
	)

	for {
		e, err := p.Event()
		if err != nil {
			return nil, err
		}
		if e.Type() == gosax.EventEOF {
			break
		}
		name, _ := gosax.Name(e.Bytes)
		switch e.Type() {
		case gosax.EventStart:
			switch string(name) {
			case "body":
				bodyNum++
				bodyName = "chapter"
				if bodyNum > 1 {
					bodyName = fmt.Sprintf("notes-%d", bodyNum-1)
				}
				updatePage()
			case "section":
				if bodyName == "chapter" && sectionDepth == 0 {
					updatePage()
				}
				sectionDepth++
			}

			if v := getAttr(e.Bytes, "id"); v != "" {
				links[`#`+v] = itemName
			}

		case gosax.EventEnd:
			if string(name) == "section" {
				sectionDepth--
			}

		}
	}
	p.chapterNum = 0
	return links, nil
}

func (p *FB2Parser) MakeEpub(wc io.WriteCloser) error {
	links, err := p.links()
	if err != nil {
		return err
	}
	err = p.Restart()
	if err != nil {
		return err
	}

	epub, err := epub2.New(wc)
	if err != nil {
		return err
	}

	defer epub.Close()

	p.parent = newTagStack()
	p.chapterNum = 0
	bodyNum := 0
	for {
		e, err := p.Event()
		if err != nil {
			return err
		}
		if e.Type() == gosax.EventEOF {
			break
		}

		name, _ := gosax.Name(e.Bytes)
		switch e.Type() {
		case gosax.EventStart:
			switch string(name) {
			case "FictionBook":
				p.parent.reset()
			case "description":
				p.parent.push("description")
				if err = p.parseDescription(epub); err != nil {
					return err
				}

			case "body":
				p.parent.push("body")
				bodyNum++
				bodyName := "chapter"
				if bodyNum > 1 {
					bodyName = fmt.Sprintf("notes-%d", bodyNum-1)
				}

				if err = p.parseBody(epub, bodyName, links); err != nil {
					return err
				}

			case "binary":
				id := getAttr(e.Bytes, "id")
				contentType := getAttr(e.Bytes, "content-type")
				content := getText(p.Reader)
				if err = epub.AddBinary(id, contentType, content); err != nil {
					return err
				}
			}
		}
	}

	if err = epub.AddTOC(); err != nil {
		return err
	}

	if err = epub.AddOPF(); err != nil {
		return err
	}
	return nil
}

func getText(r *gosax.Reader) string {
	var data []byte
	for {
		e, err := r.Event()
		if err != nil {
			return ""
		}
		if e.Type() == gosax.EventEOF {
			return ""
		}
		switch e.Type() {
		case gosax.EventText:
			data = append(data, e.Bytes...)
		case gosax.EventEnd:
			return string(bytes.TrimSpace(data))
		}
	}
}

func getAttr(b []byte, name string) string {
	_, b = gosax.Name(b)
	var attr gosax.Attribute
	var err error
	for len(b) > 0 {
		attr, b, err = gosax.NextAttribute(b)
		if err != nil {
			return ""
		}
		if i := bytes.IndexByte(attr.Key, ':'); i >= 0 {
			attr.Key = attr.Key[i+1:]
		}

		if string(attr.Key) == name {
			attr.Value, err = gosax.Unescape(attr.Value[1 : len(attr.Value)-1])
			if err != nil {
				return ""
			}
			return string(attr.Value)
		}
	}
	return ""
}
