package fb2

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/vinser/flibgolite/pkg/parser"
)

type FB2 struct {
	*TitleInfo
}

func NewFB2(rc io.ReadCloser) (*FB2, error) {
	decoder := parser.NewXmlDecoder(rc)
	fb := &FB2{}
TokenLoop:
	for {
		t, err := decoder.Token()
		if err != nil {
			return nil, err
		}

		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "title-info" {
				decoder.DecodeElement(fb, &se)
				break TokenLoop
			}
		default:
		}
	}
	return fb, nil
}

func (fb *FB2) String() string {
	return fmt.Sprint(
		"=========FB2===================\n",
		fmt.Sprintf("Authors:    %#v\n", fb.Authors),
		fmt.Sprintf("Title:      %#v\n", fb.Title),
		fmt.Sprintf("Gengre:     %#v\n", fb.Gengres),
		fmt.Sprintf("Annotation: %#v\n", fb.Annotation),
		fmt.Sprintf("Keywords:   %#v\n", fb.Keywords),
		fmt.Sprintf("Date:       %#v\n", fb.Date),
		fmt.Sprintf("Year:       %#v\n", fb.Year),
		fmt.Sprintf("Lang:       %#v\n", fb.Lang),
		fmt.Sprintf("Serie:      %#v\n", fb.Serie),
		fmt.Sprintf("CoverPage:  %#v\n", fb.CoverPage),
		"===============================\n",
	)
}

type TitleInfo struct {
	Authors    []Author   `xml:"author"`
	Title      string     `xml:"book-title"`
	Gengres    []string   `xml:"genre"`
	Annotation Annotation `xml:"annotation"`
	Keywords   string     `xml:"keywords"`
	Date       string     `xml:"date"`
	Year       string     `xml:"year"`
	Lang       string     `xml:"lang"`
	Serie      Serie      `xml:"sequence"`
	CoverPage  Image      `xml:"coverpage>image"`
}

type Author struct {
	FirstName  string `xml:"first-name"`
	MiddleName string `xml:"middle-name"`
	LastName   string `xml:"last-name"`
}

type Annotation struct {
	Text string `xml:",innerxml"`
}

type Serie struct {
	Name   string `xml:"name,attr"`
	Number int    `xml:"number,attr"`
}

type CoverPage struct {
	*Image `xml:"image"`
}

type Image struct {
	Href string `xml:"http://www.w3.org/1999/xlink href,attr"`
}

type Binary struct {
	Id          string `xml:"id,attr"`
	ContentType string `xml:"content-type,attr"`
	Content     []byte `xml:",chardata"`
}

func GetCoverPageBinary(coverLink string, rc io.ReadCloser) (*Binary, error) {
	decoder := parser.NewXmlDecoder(rc)
	b := &Binary{}
	coverLink = strings.TrimPrefix(coverLink, "#")
TokenLoop:
	for {
		t, _ := decoder.Token()
		if t == nil {
			return nil, errors.New("FB2 xml error")
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "binary" {
				for _, att := range se.Attr {
					if strings.ToLower(att.Name.Local) == "id" && att.Value == coverLink {
						decoder.DecodeElement(b, &se)
						break TokenLoop
					}
				}
			}
		default:
		}
	}
	if b == nil {
		return nil, errors.New("FB2 has no Cover Page")
	}
	return b, nil
}

func (b *Binary) String() string {
	return fmt.Sprintf(
		`CoverPage ----
  Id: %s
  Content-type: %s
================================
%#v
===========(100)================
`, b.Id, b.ContentType, b.Content[:99])
}
