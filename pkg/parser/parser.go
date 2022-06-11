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

func NewDecoder(rc io.ReadCloser) *xml.Decoder {
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

// RegExp Remove surplus spaces
var rxSpaces = regexp.MustCompile(`[ \n\r\t]+`)

func CollapceSpaces(s string) string {
	return rxSpaces.ReplaceAllString(s, ` `)
}

// RegExp Find first genre in  a string
var rxGenre = regexp.MustCompile(`[\pL\pN_]{2,}`)

func FirstGenre(s string) string {
	return rxGenre.FindString(s)
}

// RegExp Find all keywords in a string
var rxKeyword = regexp.MustCompile(`[\pL\pN]{3,}`)

func Keywords(s string) string {
	return strings.Join(rxKeyword.FindAllString(s, -1), " ")
}
