package parser

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/vinser/flibgolite/pkg/model"
	"golang.org/x/net/html"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
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

func RefineName(n, lang string) string {
	return Title(Lower(strings.TrimSpace(n), lang), lang)
}

func Title(s, lang string) string {
	return cases.Title(GetLanguageTag(lang)).String(s)
}

func Lower(s, lang string) string {
	return cases.Lower(GetLanguageTag(lang)).String(s)
}

func Upper(s, lang string) string {
	return cases.Upper(GetLanguageTag(lang)).String(s)
}

func GetSortTitle(title string, tag language.Tag) string {
	title = strings.TrimSpace(title)
	if base, _ := tag.Base(); base.String() == "en" {
		title = DropLeadingEnglishArticles(title)
	}
	title = cases.Upper(tag).String(AlphaNum(title))
	return title
}

func GetLanguageTag(lang string) language.Tag {
	return language.Make(strings.TrimSpace(lang))
}

func GetLanguage(lang string) *model.Language {
	tag := GetLanguageTag(lang)
	base, _ := tag.Base()
	langName := cases.Title(tag).String(display.Self.Name(tag))
	return &model.Language{
		Code: base.String(),
		Name: langName,
	}
}

// RegExp Remove non-alphanum unicode runes
var rxNonAlphaNum = regexp.MustCompile(`[^\pL\pN ]`)

func AlphaNum(s string) string {
	return strings.TrimSpace(rxNonAlphaNum.ReplaceAllString(s, ``))
}

// RegExp Remove surplus spaces
var rxSpaces = regexp.MustCompile(`[ \n\r\t]+`)

func CollapseSpaces(s string) string {
	return rxSpaces.ReplaceAllString(s, ` `)
}

// RegExp Find first genre in a string
var rxGenre = regexp.MustCompile(`[\pL\pN_]{2,}`)

func FirstGenre(s string) string {
	return rxGenre.FindString(s)
}

// RegExp Find first year in  a string
var rxYear = regexp.MustCompile(`[1|2][\pN]{3,3}`)

func PickYear(s string) string {
	return rxYear.FindString(s)
}

// RegExp Find all keywords in a string
var rxKeyword = regexp.MustCompile(`[\pL\pN]{3,}`)

func Keywords(s string) string {
	return strings.Join(rxKeyword.FindAllString(s, -1), ` `)
}

// RegExp Find English article at the beginning of the string
var rx1stArticle = regexp.MustCompile(`(?i)^An? |^The `)

func DropLeadingEnglishArticles(s string) string {
	return strings.TrimSpace(rx1stArticle.ReplaceAllString(s, ``))
}

func StripHTMLTags(text string) string {
	node, err := html.Parse(strings.NewReader(text))
	if err != nil {
		// If it cannot be parsed text as HTML, return the text as is.
		return text
	}

	buf := &bytes.Buffer{}
	removeHtmlTags(node, buf)

	return buf.String()
}

func removeHtmlTags(node *html.Node, buf *bytes.Buffer) {
	if node.Type == html.TextNode {
		buf.WriteString(node.Data)
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		removeHtmlTags(child, buf)
	}
}
