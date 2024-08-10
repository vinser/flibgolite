package fb2

import (
	"encoding/xml"

	"fmt"
	"strings"

	"github.com/vinser/flibgolite/pkg/conv/epub2"
)

func (p *FB2Parser) parseBody(e *epub2.EPUB, bodyName string, links map[string]string) error {
	var (
		err            error
		content        string
		currentSection string
		sectionDepth   int
		sectionNum     int
		sectionId      string

		updatePage = func() error {
			if len(content) > 0 {
				p.chapter++
				currentSection = fmt.Sprintf("%s_%d", bodyName, p.chapter)
				if err = e.AddItem(currentSection, "text", content); err != nil {
					return err
				}
				content = ""
			}
			sectionDepth = 0

			return nil
		}
	)

	updatePage()
	isNavTitle := bodyName == "chapter"

	for {
		token, err := p.Token()
		if err != nil {
			return err
		}

		switch t := token.(type) {
		case xml.StartElement:
			attrId := ""
			id := ""
			for _, a := range t.Attr {
				if a.Name.Local == "id" {
					id = a.Value
					attrId = ` id="` + id + `"`
					break
				}
			}

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
				isNavTitle = bodyName == "chapter"

			case "title":
				content += `<div class="` + t.Name.Local + `"` + attrId + `>`
				if isNavTitle { // body or section has title, add title to TOC
					if err := func() error {
						title := ""
						for {
							token, err := p.Token()
							if err != nil {
								return err
							}

							switch t := token.(type) {
							case xml.StartElement:
								switch t.Name.Local {
								case "p":
									content += "<p" + attrId + ">"
								case "empty-line":
									content += `<br />`
									title += `<br />`
								}

							case xml.CharData:
								content += string(t)
								title += string(t)

							case xml.EndElement:
								switch t.Name.Local {
								case "p":
									content += `</p>`
								case "title":
									content += `</div>`
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
											Src:   fmt.Sprintf("%s_%d.xhtml#%s", bodyName, p.chapter+1, sectionId),
											Depth: sectionDepth,
										})
									}
									return nil
								}
							}
						}
					}(); err != nil {
						return err
					}
				}
			case "epigraph", "poem", "stanza", "text-author", "subtitle", "cite":
				content += `<div class="` + t.Name.Local + `"` + attrId + `>`

			case "emphasis":
				content += "<em" + attrId + ">"

			case "p":
				content += "<p" + attrId + ">"

			case "v":
				content += `<p class="v"` + attrId + `>`

			case "empty-line":
				content += `<br />`

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
							content += `<div class="image"><img src="` + image + `"` + attrId + " /></div>"
						}
						break
					}
				}
				if isNavTitle { // body image is before title
					continue
				}
			}
		case xml.CharData:
			content += string(t)

		case xml.EndElement:
			switch t.Name.Local {
			case "section":
				sectionDepth--
				isNavTitle = false
				fallthrough

			case "epigraph", "poem", "stanza", "text-author", "title", "subtitle":
				content += `</div>`
			case "emphasis":
				content += `</em>`
			case "p", "v":
				content += `</p>`
			case "a":
				content += `</a>`
			case "body":
				return updatePage()
			}
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
			panic(fmt.Sprintf("unexpected xml.Token: %#v", t))
		}
	}
}
