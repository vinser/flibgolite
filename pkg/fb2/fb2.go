package fb2

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/vinser/flibgolite/pkg/model"

	"golang.org/x/text/cases"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/language"
)

type FB2 struct {
	*TitleInfo
}

func NewFB2(rc io.ReadCloser) (*FB2, error) {
	decoder := xml.NewDecoder(rc)
	decoder.CharsetReader = charsetReader
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
	decoder := xml.NewDecoder(rc)
	decoder.CharsetReader = charsetReader
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

func charsetReader(charset string, input io.Reader) (io.Reader, error) {
	switch strings.ToLower(charset) {
	case "windows-1251":
		return charmap.Windows1251.NewDecoder().Reader(input), nil
	case "windows-1252":
		return charmap.Windows1252.NewDecoder().Reader(input), nil
	default:
		return nil, fmt.Errorf("unknown charset: %s", charset)
	}
}

func (fb *FB2) GetFormat() string {
	return "fb2"
}

func (fb *FB2) GetTitle() string {
	return strings.Trim(fb.Title, "\n\t ")
}

func (fb *FB2) GetSort() string {
	return strings.ToUpper(strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(strings.Trim(fb.Title, "\n\t "), "An "), "A "), "The "))
}

func (fb *FB2) GetYear() string {
	year := fb.Year
	if year == "" {
		year = fb.Date
	}
	rYear := []rune(year)
	if len(rYear) > 4 {
		rYear = rYear[len(rYear)-4:]
	}
	return strings.Trim(string(rYear), "\n\t ")
}

func (fb *FB2) GetPlot() string {
	return fb.Annotation.Text
}

func (fb *FB2) GetCover() string {
	return strings.TrimPrefix(fb.CoverPage.Href, "#")
}

func (fb *FB2) GetLanguage() *model.Language {
	base, _ := fb.getLanguageTag().Base()
	return &model.Language{Code: base.String()}
}

func (fb *FB2) GetAuthors() []*model.Author {
	authors := make([]*model.Author, 0, len(fb.Authors))
	if len(fb.Authors) == 1 {
		aLN := strings.Split(fb.Authors[0].LastName, ",")
		if len(aLN) > 1 {
			a := "Авторский коллектив"
			if fb.Lang != "ru" {
				a = "Writing team"
			}
			authors = append(authors, &model.Author{
				Name: a,
				Sort: strings.ToUpper(a),
			})
			return authors
		}
	}
	for _, a := range fb.Authors {
		author := &model.Author{}
		f := fb.title(fb.lower(strings.Trim(a.FirstName, "\n\t ")))
		m := fb.title(fb.lower(strings.Trim(a.MiddleName, "\n\t ")))
		l := fb.title(fb.lower(strings.Trim(a.LastName, "\n\t ")))
		author.Name = strings.Trim(strings.ReplaceAll(fmt.Sprint(f, " ", m, " ", l), "  ", " "), " ")
		author.Sort = strings.Trim(strings.ReplaceAll(fmt.Sprint(l, " ", f, " ", m), "  ", " "), " ")
		authors = append(authors, author)
	}
	return authors
}

func (fb *FB2) GetGenres() []string {
	return fb.Gengres
}

func (fb *FB2) GetSerie() *model.Serie {
	return &model.Serie{Name: fb.title(fb.Serie.Name)}
}

func (fb *FB2) GetSerieNumber() int {
	return fb.Serie.Number
}

func (fb *FB2) getLanguageTag() language.Tag {
	code := strings.Trim(fb.Lang, "\n\t ")
	if strings.TrimSpace(code) == "uk" { // patch old "uk" for Ukrainian to morden "ua"
		code = "au"
	}
	return language.Make(code)
}

func (fb *FB2) title(s string) string {
	c := cases.Title(fb.getLanguageTag())
	return c.String(s)
}

func (fb *FB2) lower(s string) string {
	c := cases.Lower(fb.getLanguageTag())
	return c.String(s)
}
