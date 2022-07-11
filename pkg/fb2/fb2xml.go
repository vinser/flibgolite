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
	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/vinser/flibgolite/pkg/model"
	"github.com/vinser/flibgolite/pkg/parser"
)

type FB2 struct {
	FictionBook xml.Name `xml:"FictionBook"` // Root element
	Description struct {
		TitleInfo struct { // Generic information about the book
			Authors []struct { // Author(s) of this boo
				FirstName  string `xml:"first-name"`
				MiddleName string `xml:"middle-name"`
				LastName   string `xml:"last-name"`
			} `xml:"author"`
			BookTitle  string   `xml:"book-title"` // Book title
			Genres     []string `xml:"genre"`      // Genre of this book
			Annotation struct { // Annotation for this book
				P []string `xml:"p"`
			} `xml:"annotation"`
			Keywords string     `xml:"keywords"` // Any keywords for this book, intended for use in search engines
			Date     string     `xml:"date"`     // Date this book was written, can be not exact, e.g. 1863-1867.
			Lang     string     `xml:"lang"`     // Book language
			Series   []struct { // Any sequences this book might be part of
				Name   string `xml:"name,attr"`
				Number int    `xml:"number,attr"`
			} `xml:"sequence"`
			CoverPage struct { // Any coverpage items, currently only images
				Image struct {
					Href string `xml:"href,attr"`
				} `xml:"image"`
			} `xml:"coverpage"`
		} `xml:"title-info"`

		PublishInfo struct { // Information about some paper/outher published document, that was used as a source of this xml document
			Publisher string     `xml:"publisher"` // Original (paper) book publisher
			City      string     `xml:"city"`      // City where the original (paper) book was published
			Year      int        `xml:"year"`      // Year of the original (paper) publication
			ISBN      string     `xml:"isbn"`      // ISBN
			Series    []struct { // Any sequences this book might be part of
				Name   string `xml:"name,attr"`
				Number int    `xml:"number,attr"`
			} `xml:"sequence"`
		} `xml:"publish-info"`
	} `xml:"description"`
}

func NewFB2(rc io.ReadCloser) (*FB2, error) {
	decoder := parser.NewDecoder(rc)
	fb := &FB2{}
	// decoder.Decode(&fb)
TokenLoop:
	for {
		t, err := decoder.Token()
		if t == nil {
			return nil, err
		}

		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "FictionBook" {
				decoder.DecodeElement(fb, &se)
				break TokenLoop
			}
		default:
		}
	}
	return fb, nil
}

func (fb *FB2) String() string {
	return "" + fmt.Sprint(
		"=========FB2===================\n",
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
		fmt.Sprintf("Publisher:  %#v\n", fb.Description.PublishInfo.Publisher),
		fmt.Sprintf("City:       %#v\n", fb.Description.PublishInfo.City),
		fmt.Sprintf("Year:       %#v\n", fb.Description.PublishInfo.Year),
		fmt.Sprintf("ISBN:       %#v\n", fb.Description.PublishInfo.ISBN),
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
	if book.Archive == "" {
		rc, _ = os.Open(path.Join(stock, book.File))
	} else {
		zr, _ := zip.OpenReader(path.Join(stock, book.Archive))
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
