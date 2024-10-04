package fb2

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/orisano/gosax"
	"github.com/vinser/u8xml"
)

type FB2 struct {
	// FictionBook xml.Name `xml:"FictionBook"` // Root element
	Description struct {
		TitleInfo   TitleInfo
		PublishInfo PublishInfo
	}
}
type TitleInfo struct { // Generic information about a book
	Authors    []Author   // Author(s) of a book
	BookTitle  string     // Book title
	Genres     []string   // Genre of a book
	Annotation Annotation // Annotation of a book
	Keywords   string     // Any keywords of a book, intended for use in search engines
	Date       string     // Date a book was written, can be not exact, e.g. 1863-1867.
	Lang       string     // Book language
	Series     []Serie    // Any sequences a book might be part of
	CoverPage  CoverPage  // Any coverpage items, currently only images
}
type PublishInfo struct { // Information about some paper/outher published document, that was used as a source of this xml document
	Year   int     // Year of the original (paper) publication
	Series []Serie // Any sequences a book might be part of
}

type Annotation struct { // Annotation of a book
	P []string
}

type Author struct { // Author of a book
	FirstName  string
	MiddleName string
	LastName   string
}

type Serie struct { // Any sequences this book might be part of
	Name   string
	Number int
}

type CoverPage struct { // Any coverpage items, currently only images
	Image struct {
		Href string
	}
}

var (
	ErrNoElement     = errors.New("no element")
	ErrUnexpectedEOF = errors.New("unexpected EOF")
)

func ParseFB2Description(rc io.ReadCloser) (*FB2, error) {
	u8r, err := u8xml.NewReader(rc)
	if err != nil {
		return nil, err
	}
	fb := &FB2{}
	r := gosax.NewReader(u8r)
	r.EmitSelfClosingTag = true
	for {
		e, err := r.Event()
		if err != nil {
			return nil, err
		}
		if e.Type() == gosax.EventEOF {
			return nil, ErrUnexpectedEOF
		}

		name, _ := gosax.Name(e.Bytes)
		switch e.Type() {
		case gosax.EventStart:
			switch string(name) {
			case "title-info":
				titleInfo, err := parseTitleInfo(r)
				if err == nil {
					fb.Description.TitleInfo = titleInfo
				}
			case "publish-info":
				publishInfo, err := parsePublishInfo(r)
				if err == nil {
					fb.Description.PublishInfo = publishInfo
				}
			}
		case gosax.EventEnd:
			if string(name) == "description" {
				return fb, nil
			}
		}
	}
}

func parseTitleInfo(r *gosax.Reader) (TitleInfo, error) {
	titleInfo := TitleInfo{}
	for {
		e, err := r.Event()
		if err != nil {
			return titleInfo, err
		}
		if e.Type() == gosax.EventEOF {
			return titleInfo, ErrUnexpectedEOF
		}
		name, _ := gosax.Name(e.Bytes)
		switch e.Type() {
		case gosax.EventStart:
			switch string(name) {
			case "author":
				author, err := parseAuthor(r)
				if err == nil {
					titleInfo.Authors = append(titleInfo.Authors, author)
				}
			case "book-title":
				titleInfo.BookTitle = getText(r)
			case "genre":
				genre := getText(r)
				if genre != "" {
					titleInfo.Genres = append(titleInfo.Genres, genre)
				}
			case "annotation":
				annotation, err := parseAnnotation(r)
				if err == nil {
					titleInfo.Annotation = annotation
				}
			case "keywords":
				titleInfo.Keywords = getText(r)
			case "date":
				titleInfo.Date = getText(r)
			case "lang":
				titleInfo.Lang = getText(r)
			case "sequence":
				serie, err := parseSerie(e.Bytes)
				if err == nil {
					titleInfo.Series = append(titleInfo.Series, serie)
				}
			case "coverpage":
				cover, err := parseCoverPage(r)
				if err == nil {
					titleInfo.CoverPage = cover
				}
			}

		case gosax.EventEnd:
			if string(name) == "title-info" {
				return titleInfo, nil
			}
		}
	}
}

func parsePublishInfo(r *gosax.Reader) (PublishInfo, error) {
	publishInfo := PublishInfo{}
	for {
		e, err := r.Event()
		if err != nil {
			return publishInfo, err
		}
		if e.Type() == gosax.EventEOF {
			return publishInfo, ErrUnexpectedEOF
		}
		name, _ := gosax.Name(e.Bytes)
		switch e.Type() {
		case gosax.EventStart:
			switch string(name) {
			case "year":
				publishInfo.Year, _ = strconv.Atoi(getText(r))
			case "series":
				serie, err := parseSerie(e.Bytes)
				if err == nil {
					publishInfo.Series = append(publishInfo.Series, serie)
				}
			}
		case gosax.EventEnd:
			if string(name) == "publish-info" {
				return publishInfo, nil
			}
		}
	}
}

func parseAuthor(r *gosax.Reader) (Author, error) {
	author := Author{}
	for {
		e, err := r.Event()
		if err != nil {
			return author, err
		}
		if e.Type() == gosax.EventEOF {
			return author, ErrUnexpectedEOF
		}
		name, _ := gosax.Name(e.Bytes)
		switch e.Type() {
		case gosax.EventStart:
			switch string(name) {
			case "first-name":
				author.FirstName = getText(r)
			case "middle-name":
				author.MiddleName = getText(r)
			case "last-name":
				author.LastName = getText(r)
			}
		case gosax.EventEnd:
			switch string(name) {
			case "author":
				return author, nil
			}
		}
	}
}

func parseAnnotation(r *gosax.Reader) (Annotation, error) {
	annotation := Annotation{}
	for {
		e, err := r.Event()
		if err != nil {
			return annotation, err
		}
		if e.Type() == gosax.EventEOF {
			return annotation, ErrUnexpectedEOF
		}
		name, _ := gosax.Name(e.Bytes)
		switch e.Type() {
		case gosax.EventStart:
			switch string(name) {
			case "p":
				if bytes.IndexByte(e.Bytes, '/') >= 0 {
					annotation.P = append(annotation.P, "")
				} else {
					if p := getText(r); p != "" {
						annotation.P = append(annotation.P, p)
					}
				}
			case "empty-line":
				annotation.P = append(annotation.P, "\n")
			}
			// TODO: parse other tags
		case gosax.EventEnd:
			switch string(name) {
			case "annotation":
				return annotation, nil
			}
		}
	}
}

func parseSerie(b []byte) (Serie, error) {
	name := getAttr(b, "name")
	if name != "" {
		number, err := strconv.Atoi(getAttr(b, "number"))
		if err != nil {
			return Serie{}, err
		}
		return Serie{
			Name:   name,
			Number: number,
		}, nil
	}
	return Serie{}, ErrNoElement
}

func parseCoverPage(r *gosax.Reader) (CoverPage, error) {
	coverPage := CoverPage{}
	for {
		e, err := r.Event()
		if err != nil {
			return coverPage, err
		}
		if e.Type() == gosax.EventEOF {
			return coverPage, ErrUnexpectedEOF
		}
		name, _ := gosax.Name(e.Bytes)
		switch e.Type() {
		case gosax.EventStart:
			if string(name) == "image" {
				href := getAttr(e.Bytes, "href")
				if href != "" {
					coverPage.Image.Href = href
				}
			}
		case gosax.EventEnd:
			switch string(name) {
			case "coverpage":
				return coverPage, nil
			}
		}
	}
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

// =================================
func (fb *FB2) String() string {
	return "" + fmt.Sprint(
		"\n=========FB2===================\n",
		"---------TitleInfo-------------\n",
		fmt.Sprintf("Authors:    %#v\n", fb.Description.TitleInfo.Authors),
		fmt.Sprintf("BookTitle:  %#v\n", fb.Description.TitleInfo.BookTitle),
		fmt.Sprintf("Genres:    %#v\n", fb.Description.TitleInfo.Genres),
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
