// Made on scheme https://github.com/gribuser/fb2 licenced by Dmitry Gribov

// Copyright (c) 2004, Dmitry Gribov
// All rights reserved.

// Redistribution and use in source and binary forms, with or without modification,
// are permitted provided that the following conditions are met:

//     * Redistributions of source code must retain the above copyright notice, this list
//     of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above copyright notice, this
//     list of conditions and the following disclaimer in the documentation and/or other
//     materials provided with the distribution.

// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY
// EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES
// OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT
// SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED
// TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR
// BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN
// ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
// SUCH DAMAGE.

package fb2

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"image"
	"io"
	"os"
	"path"
	"strconv"
	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/vinser/flibgolite/pkg/model"
	"github.com/vinser/flibgolite/pkg/parser"
	"golang.org/x/net/html/charset"
)

type FB2 struct {
	// FictionBook xml.Name `xml:"FictionBook"` // Root element
	Description struct {
		TitleInfo   TitleInfo   `xml:"title-info"`
		PublishInfo PublishInfo `xml:"publish-info"`
	} `xml:"description"`
}
type TitleInfo struct { // Generic information about a book
	Authors    []Author   `xml:"author"`     // Author(s) of a book
	BookTitle  string     `xml:"book-title"` // Book title
	Genres     []string   `xml:"genre"`      // Genre of a book
	Annotation Annotation `xml:"annotation"` // Annotation of a book
	Keywords   string     `xml:"keywords"`   // Any keywords of a book, intended for use in search engines
	Date       string     `xml:"date"`       // Date a book was written, can be not exact, e.g. 1863-1867.
	Lang       string     `xml:"lang"`       // Book language
	Series     []Serie    `xml:"sequence"`   // Any sequences a book might be part of
	CoverPage  CoverPage  `xml:"coverpage"`  // Any coverpage items, currently only images
}
type PublishInfo struct { // Information about some paper/outher published document, that was used as a source of this xml document
	Year   int     `xml:"year"`     // Year of the original (paper) publication
	Series []Serie `xml:"sequence"` // Any sequences a book might be part of
}

type Annotation struct { // Annotation of a book
	P []string `xml:"p"`
}

type Author struct { // Author of a book
	FirstName  string `xml:"first-name"`
	MiddleName string `xml:"middle-name"`
	LastName   string `xml:"last-name"`
}

type Serie struct { // Any sequences this book might be part of
	Name   string `xml:"name,attr"`
	Number int    `xml:"number,attr"`
}

type CoverPage struct { // Any coverpage items, currently only images
	Image struct {
		Href string `xml:"href,attr"`
	} `xml:"image"`
}

func NewDecoder(rc io.ReadCloser) *xml.Decoder {
	decoder := xml.NewDecoder(rc)
	decoder.Strict = false
	decoder.CharsetReader = charset.NewReaderLabel
	return decoder
}

var (
	ErrNoElement = errors.New("no element")
)

func ParseFB2(rc io.ReadCloser) (*FB2, error) {
	d := NewDecoder(rc)
	fb := &FB2{}
TokenLoop:
	for {
		tok, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "FictionBook":
			case "description":
			case "title-info":
				titleInfo, err := parseTitleInfo(d)
				if err == nil {
					fb.Description.TitleInfo = titleInfo
				}
			case "publish-info":
				publishInfo, err := parsePublishInfo(d)
				if err == nil {
					fb.Description.PublishInfo = publishInfo
				}
			default:
				d.Skip()
			}
		case xml.EndElement:
			switch t.Name.Local {
			case "description":
				break TokenLoop
			}
		}
	}
	return fb, nil
}

func parseTitleInfo(d *xml.Decoder) (TitleInfo, error) {
	titleInfo := TitleInfo{}
	for {
		tok, err := d.Token()
		if err != nil {
			return titleInfo, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "author":
				author, err := parseAuthor(d)
				if err == nil {
					titleInfo.Authors = append(titleInfo.Authors, author)
				}
			case "book-title":
				titleInfo.BookTitle = getValue(d)
			case "genre":
				genre := getValue(d)
				if genre != "" {
					titleInfo.Genres = append(titleInfo.Genres, genre)
				}
			case "annotation":
				annotation, err := parseAnnotation(d)
				if err == nil {
					titleInfo.Annotation = annotation
				}
			case "keywords":
				titleInfo.Keywords = getValue(d)
			case "date":
				titleInfo.Date = getValue(d)
			case "lang":
				titleInfo.Lang = getValue(d)
			case "sequence":
				serie, err := parseSerie(t)
				if err == nil {
					titleInfo.Series = append(titleInfo.Series, serie)
				}
			case "coverpage":
				cover, err := parseCoverPage(d)
				if err == nil {
					titleInfo.CoverPage = cover
				}
			default:
				d.Skip()
			}
		case xml.EndElement:
			switch t.Name.Local {
			case "title-info":
				return titleInfo, nil
			}
		}
	}
}

func parsePublishInfo(d *xml.Decoder) (PublishInfo, error) {
	publishInfo := PublishInfo{}
	for {
		tok, err := d.Token()
		if err != nil {
			return publishInfo, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "year":
				publishInfo.Year, _ = strconv.Atoi(getValue(d))
			case "series":
				serie, err := parseSerie(t)
				if err == nil {
					publishInfo.Series = append(publishInfo.Series, serie)
				}
			default:
				d.Skip()
			}
		case xml.EndElement:
			switch t.Name.Local {
			case "publish-info":
				return publishInfo, nil
			}
		}
	}
}

func parseAuthor(d *xml.Decoder) (Author, error) {
	author := Author{}
	for {
		tok, err := d.Token()
		if err != nil {
			return author, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "first-name":
				author.FirstName = getValue(d)
			case "middle-name":
				author.MiddleName = getValue(d)
			case "last-name":
				author.LastName = getValue(d)
			default:
				d.Skip()
			}
		case xml.EndElement:
			switch t.Name.Local {
			case "author":
				return author, nil
			}
		}
	}
}

func parseAnnotation(d *xml.Decoder) (Annotation, error) {
	annotation := Annotation{}
	for {
		tok, err := d.Token()
		if err != nil {
			return annotation, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "p":
				p := getValue(d)
				if p != "" {
					annotation.P = append(annotation.P, p)
				}
			default:
				d.Skip()
			}
		case xml.EndElement:
			switch t.Name.Local {
			case "annotation":
				return annotation, nil
			}
		}
	}
}

func parseSerie(token xml.StartElement) (Serie, error) {
	name := getAttr(token, "name")
	if name != "" {
		number, _ := strconv.Atoi(getAttr(token, "number"))
		return Serie{
			Name:   name,
			Number: number,
		}, nil
	}
	return Serie{}, ErrNoElement
}

func parseCoverPage(d *xml.Decoder) (CoverPage, error) {
	coverPage := CoverPage{}
	for {
		tok, err := d.Token()
		if err != nil {
			return coverPage, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			if t.Name.Local == "image" {
				href := getAttr(t, "href")
				if href != "" {
					coverPage.Image.Href = href
				}
			}
		case xml.EndElement:
			if t.Name.Local == "coverpage" {
				return coverPage, nil
			}
		}
	}
}

func getValue(d xml.TokenReader) string {
	var data []byte
	for {
		tok, err := d.Token()
		if err != nil {
			break
		}
		switch t := tok.(type) {
		case xml.CharData:
			data = append(data, t...)
		case xml.EndElement:
			return string(bytes.TrimSpace(data))
		}
	}
	return ""
}

func getAttr(e xml.StartElement, name string) string {
	for _, a := range e.Attr {
		if a.Name.Local == name {
			return a.Value
		}
	}
	return ""
}
func (fb *FB2) String() string {
	return "" + fmt.Sprint(
		"\n=========FB2===================\n",
		"---------TitleInfo-------------\n",
		fmt.Sprintf("Authors:    %#v\n", fb.Description.TitleInfo.Authors),
		fmt.Sprintf("BookTitle:  %#v\n", fb.Description.TitleInfo.BookTitle),
		fmt.Sprintf("Gengres:    %#v\n", fb.Description.TitleInfo.Genres),
		fmt.Sprintf("Annotation: %#v\n", fb.Description.TitleInfo.Annotation),
		fmt.Sprintf("Keywords:   %#v\n", fb.Description.TitleInfo.Keywords),
		fmt.Sprintf("Date:       %#v\n", fb.Description.TitleInfo.Date),
		fmt.Sprintf("Lang:       %#v\n", fb.Description.TitleInfo.Lang),
		fmt.Sprintf("Series:     %#v\n", fb.Description.TitleInfo.Series),
		fmt.Sprintf("CoverPage:  %#v\n", fb.Description.TitleInfo.CoverPage),
		"---------PublishInfo-----------\n",
		fmt.Sprintf("Year:       %#v\n", fb.Description.PublishInfo.Year),
		fmt.Sprintf("Series:     %#v\n", fb.Description.PublishInfo.Series),
		"===============================\n",
	)
}

// Any binary data that is required for the presentation of this book in base64 format.
// Currently only images are used.
type Binary struct {
	ID          string `xml:"id,attr"`
	ContentType string `xml:"content-type,attr"`
	Content     []byte `xml:",chardata"`
}

func getCoverPageBinary(coverLink string, rc io.ReadCloser) (*Binary, error) {
	decoder := parser.NewDecoder(rc)
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
	return b, nil
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
	if book.Archive.Name == "" {
		rc, _ = os.Open(path.Join(stock, book.File))
	} else {
		zr, _ := zip.OpenReader(path.Join(stock, book.Archive.Name))
		defer zr.Close()
		for _, file := range zr.File {
			if file.Name == book.File {
				rc, _ = file.Open()
				break
			}
		}
	}
	defer rc.Close()
	b, err := getCoverPageBinary(book.Cover, rc)
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
