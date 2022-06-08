package parser

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/vinser/flibgolite/pkg/model"
	"golang.org/x/text/cases"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/language"
)

type Parser interface {
	GetFormat() string
	GetTitle() string
	GetSort() string
	GetYear() string
	GetPlot() string
	GetCover() string
	GetLanguage() *model.Language
	GetAuthors() []*model.Author
	GetGenres() []string
	GetKeywords() string
	GetSerie() *model.Serie
	GetSerieNumber() int
}

func NewXmlDecoder(rc io.ReadCloser) *xml.Decoder {
	decoder := xml.NewDecoder(rc)
	decoder.CharsetReader = charsetReader
	return decoder
}

func charsetReader(charset string, input io.Reader) (io.Reader, error) {
	switch strings.ToLower(charset) {
	case "windows-1251":
		return charmap.Windows1251.NewDecoder().Reader(input), nil
	case "windows-1252":
		return charmap.Windows1252.NewDecoder().Reader(input), nil
	case "iso-8859-5":
		return charmap.ISO8859_5.NewDecoder().Reader(input), nil
	default:
		return nil, fmt.Errorf("unknown charset: %s", charset)
	}
}

func RefineName(n, lang string) string {
	return Title(Lower(strings.TrimSpace(n), lang), lang)
}

func Title(s, lang string) string {
	return cases.Title(GetLanguageTag(lang)).String(s)
}

func Lower(s, lang string) string {
	return cases.Lower(GetLanguageTag(lang)).String(s)
}

func GetLanguageTag(lang string) language.Tag {
	return language.Make(strings.TrimSpace(lang))
}
