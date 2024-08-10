package epub2

import (
	"encoding/base64"
	"strconv"
)

// AddMetadataSubject adds subject to metadata
func (e *EPUB) AddMetadataSubject(subj string) {
	e.Metadata += "<dc:subject>" + subj + "</dc:subject>\n"
}

// AddMetadataAuthor adds author to metadata
func (e *EPUB) AddMetadataAuthor(name, sort string) {
	e.Metadata += `<dc:creator opf:file-as="` + sort + `" opf:role="aut" xmlns:opf="http://www.idpf.org/2007/opf">` + name + `</dc:creator>` + "\n"
}

// AddMetadataDescription adds description to metadata
func (e *EPUB) AddMetadataDescription(desc string) {
	e.Metadata += "<dc:description>" + desc + "</dc:description>\n"
}

// AddMetadataTitle adds title to metadata
func (e *EPUB) AddMetadataTitle(title string) {
	e.Title = title
	e.Metadata += "<dc:title>" + title + "</dc:title>\n"
}

// AddMetadataLanguage adds language to metadata
func (e *EPUB) AddMetadataLanguage(lang string) {
	e.Lang = lang
	e.Metadata += `<dc:language xsi:type="dcterms:RFC3066">` + lang + `</dc:language>` + "\n"
}

// AddMetadataCover add cover to metadata
func (e *EPUB) AddMetadataCover(imageName string) {
	e.Metadata += `<meta name="cover" content="` + imageName + `" />` + "\n"
}

func (e *EPUB) AddItem(itemName, guideType, content string) error {
	e.Manifest += `<item id="` + itemName + `" href="` + itemName + `.xhtml" media-type="application/xhtml+xml" />` + "\n"
	e.Spine += `<itemref idref="` + itemName + `" />` + "\n"

	data := struct {
		Title   string
		Content string
		Type    string
	}{
		Title:   e.Title,
		Content: content,
		Type:    guideType,
	}
	return e.execTemplate("OEBPS/"+itemName+".xhtml", "page.tmpl", data)
}

func (e *EPUB) AddBinary(id, contentType, base64Content string) error {
	e.Manifest += `<item id="` + id + `" href="` + id + `" media-type="` + contentType + `" />` + "\n"
	data, err := base64.StdEncoding.DecodeString(base64Content)
	if err != nil {
		return err
	}

	f, err := e.zw.Create("OEBPS/" + id)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	return err
}

func (e *EPUB) AddOPF() error {
	return e.execTemplate("OEBPS/content.opf", "content.tmpl", e)
}

func (e *EPUB) AddTOC() error {
	toc := ""
	prevDepth := 1
	for i, t := range e.Toc {
		// fmt.Println(t)
		switch {
		case i == 0 || t.Depth > prevDepth:
			prevDepth = t.Depth
		case t.Depth == prevDepth:
			toc += `
</navPoint>
`
		case t.Depth < prevDepth:
			toc += `
</navPoint>
</navPoint>
`
			prevDepth = t.Depth
		}
		toc += `
<navPoint playOrder="` + strconv.Itoa(t.Order) + `" id="` + t.Id + `">
	<navLabel>
		<text>` + t.Text + `</text>
	</navLabel>
	<content src="` + t.Src + `" />
`
	}

	for i := 1; i <= prevDepth; i++ {
		toc += `
</navPoint>
`
	}

	data := struct {
		Lang  string
		Title string
		Toc   string
	}{
		Lang:  e.Lang,
		Title: e.Title,
		Toc:   toc,
	}
	return e.execTemplate("OEBPS/toc.ncx", "toc.tmpl", data)
}
