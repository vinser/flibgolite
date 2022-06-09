package parser

import (
	"encoding/xml"
	"io"
	"regexp"
	"strings"

	"github.com/vinser/flibgolite/pkg/model"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/cases"
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
	decoder.Strict = false
	decoder.CharsetReader = charset.NewReaderLabel
	return decoder
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

var rxSpaces = regexp.MustCompile(`[ \n\r\t]+`)
var rxKeywords = regexp.MustCompile(`[\pL\pN]{3,}`)

// RegExp Remove surplus spaces
func CollapceSpaces(s string) string {
	return rxSpaces.ReplaceAllString(s, ` `)
}

// RegExp Split string to keywords slice
func ListKeywords(s string) []string {
	return rxKeywords.FindAllString(s, -1)
}
