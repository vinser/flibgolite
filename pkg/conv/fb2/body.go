package fb2

import (
	"bytes"

	"fmt"
	"strings"

	"github.com/orisano/gosax"
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
	defer p.parent.pop()

	updatePage()
	findNavTitle := bodyName == "chapter"
	insideNavTitle := false
	title := ""
	for {
		ev, err := p.Event()
		if err != nil {
			return err
		}
		name, _ := gosax.Name(ev.Bytes)
		switch ev.Type() {
		case gosax.EventStart:
			id, attrId := getAttrId(ev.Bytes)
			switch string(name) {
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
				content += `<div class="section" id="` + sectionId + `">`
				sectionDepth++
				findNavTitle = bodyName == "chapter"

			case "title":
				content += `<div class="title"` + attrId + `>`

				if findNavTitle { // body or section has title, add this title to TOC
					insideNavTitle = true
					title = ""
				}
			case "epigraph", "poem", "stanza", "text-author", "subtitle", "cite":
				content += `<div class="` + string(name) + `"` + attrId + `>`

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
				value := getAttr(ev.Bytes, "href")
				if page, ok := links[value]; ok {
					link = page + ".xhtml" + value
				} else {
					link = value
				}
				content += `<a href="` + link + `"` + attrId + ">"

			case "image":
				image := strings.TrimLeft(getAttr(ev.Bytes, "href"), "#")
				if len(image) > 0 {
					if p.parent.top() == "p" {
						content += fmt.Sprintf("<img src=\"%s\" %s alt=\"%s\" />", image, attrId, image)
					} else {
						content += fmt.Sprintf("<div class=\"image\"><img src=\"%s\" %s alt=\"%s\" /></div>", image, attrId, image)
					}
				}
				if findNavTitle { // body image is before title
					continue
				}
				p.parent.pop()
			}
			p.parent.push(string(name))

		case gosax.EventText:
			s := fixCharData(string(bytes.TrimSpace(ev.Bytes)))
			content += s
			if insideNavTitle && p.parent.top() == "p" {
				title += s + " "
			}

		case gosax.EventEnd:
			switch string(name) {
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
		}
	}
}

func getAttrId(e []byte) (value, attrId string) {
	value = getAttr(e, "id")
	if value != "" {
		attrId = ` id="` + value + `" `
	}
	return value, attrId
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
