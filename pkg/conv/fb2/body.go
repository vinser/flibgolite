package fb2

import (
	"encoding/xml"

	"fmt"
	"strings"

	"github.com/vinser/flibgolite/pkg/conv/epub2"
)

func (p *FB2Parser) parseBody(e *epub2.EPUB, bodyName string, links map[string]string) error {
	var (
		err          error
		content      string
		itemName     string
		sectionDepth int
		sectionNum   int = 1
		sectionId    string

		updatePage = func() error {
			content = strings.TrimSpace(content)
			if len(content) > 0 {
				if bodyName == "chapter" {
					p.chapterNum++
					itemName = fmt.Sprintf("%s_%d", bodyName, p.chapterNum)
				} else {
					itemName = bodyName
				}
				if err = e.AddItem(itemName, "text", content); err != nil {
					return err
				}
				content = ""
			}
			sectionDepth = 0

			return nil
		}
	)

	updatePage()
	findNavTitle := bodyName == "chapter"
	insideNavTitle := false
	title := ""
	for {
		token, err := p.Token()
		if err != nil {
			return err
		}

		switch t := token.(type) {
		case xml.StartElement:
			id, attrId := getAttrId(t)
			switch t.Name.Local {
			case "section":
				sectionNum++
				if bodyName == "chapter" && sectionDepth == 0 {
					if err = updatePage(); err != nil {
						return err
					}
				}
				if id == "" {
					sectionId = fmt.Sprintf(`section_%d`, sectionNum)
				} else {
					sectionId = id
				}
				content += `<div class="` + t.Name.Local + `" id="` + sectionId + `">`
				sectionDepth++
				findNavTitle = bodyName == "chapter"

			case "title":
				content += `<div class="` + t.Name.Local + `"` + attrId + `>`

				if findNavTitle { // body or section has title, add this title to TOC
					insideNavTitle = true
					title = ""
				}
			case "epigraph", "poem", "stanza", "text-author", "subtitle", "cite":
				content += `<div class="` + t.Name.Local + `"` + attrId + `>`

			case "emphasis":
				content += "<em" + attrId + ">"

			case "p":
				content += "<p" + attrId + ">"

			case "empty-line":
				content += `<br />`
				if insideNavTitle {
					title += "\n"
				}

			case "v":
				content += `<p class="v"` + attrId + `>`

			case "a":
				var link string
				for _, a := range t.Attr {
					if a.Name.Local == "href" {
						if page, ok := links[a.Value]; ok {
							link = page + ".xhtml" + a.Value

						} else {
							link = a.Value
						}
						break
					}
				}
				content += `<a href="` + link + `"` + attrId + ">"

			case "image":
				for _, a := range t.Attr {
					if a.Name.Local == "href" {
						image := strings.TrimLeft(a.Value, "#")
						if len(image) > 0 {
							if p.parent.top() == "p" {
								content += fmt.Sprintf("<img src=\"%s\" %s alt=\"%s\" />", image, attrId, image)
							} else {
								content += fmt.Sprintf("<div class=\"image\"><img src=\"%s\" %s alt=\"%s\" /></div>", image, attrId, image)
							}
						}
						break
					}
				}
				if findNavTitle { // body image is before title
					continue
				}
				p.parent.pop()
			}
			p.parent.push(t.Name.Local)

		case xml.CharData:
			s := fixCharData(string(t))
			content += s
			if insideNavTitle && p.parent.top() == "p" {
				title += s + " "
			}

		case xml.EndElement:
			switch t.Name.Local {
			case "section":
				sectionDepth--
				findNavTitle = false
				content += `</div>`
			case "title":
				if insideNavTitle { // body or section has title, add title to TOC
					if sectionId == "" {
						e.Toc = append(e.Toc, epub2.TOC{
							Id:    "root",
							Order: sectionNum,
							Text:  title,
							Src:   "chapter_1.xhtml",
							Depth: 1,
						})
					} else {
						e.Toc = append(e.Toc, epub2.TOC{
							Id:    sectionId,
							Order: sectionNum,
							Text:  title,
							Src:   fmt.Sprintf("%s_%d.xhtml#%s", bodyName, p.chapterNum+1, sectionId),
							Depth: sectionDepth,
						})
					}
				}
				insideNavTitle = false
				findNavTitle = false
				content += `</div>`
			case "epigraph", "poem", "stanza", "text-author", "subtitle", "cite":
				content += `</div>`
			case "emphasis":
				content += `</em>`
			case "p":
				content += `</p>`
			case "v":
				content += `</p>`
			case "a":
				content += `</a>`
			case "body":
				return updatePage()
			}
			p.parent.pop()
		case xml.Attr:
		case xml.Comment:
		case xml.Decoder:
		case xml.Directive:
		case xml.Encoder:
		case xml.Name:
		case xml.ProcInst:
		case xml.SyntaxError:
		case xml.TagPathError:
		case xml.UnmarshalError:
		case xml.UnsupportedTypeError:
		default:
			return fmt.Errorf("unexpected xml.Token: %#v", t)
		}
	}
}

func getAttrId(e xml.StartElement) (id, attrId string) {
	id = getAttrValue(e, "id")
	if id != "" {
		attrId = ` id="` + id + `" `
	}
	return id, attrId
}

func getAttrValue(e xml.StartElement, name string) string {
	for _, a := range e.Attr {
		if a.Name.Local == name {
			return a.Value
		}
	}
	return ""
}

func fixCharData(s string) string {
	s = strings.ReplaceAll(s, "& ", "&amp; ")
	transform := func(r rune) rune {
		switch r {
		case '<', '>':
			return -1
		default:
			return r
		}
	}
	return strings.Trim(strings.Map(transform, s), " \t")
}
