package parser

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/vinser/flibgolite/pkg/model"
	"golang.org/x/net/html"
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

func GetLanguageTag(lang string) language.Tag {
	return language.Make(strings.TrimSpace(lang))
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

func Drop1stArticle(s string) string {
	return strings.TrimSpace(rx1stArticle.ReplaceAllString(s, ``))
}

// RegExp Clear full name
var rxDropSpaces = regexp.MustCompile(` `)
var rxFindNames = regexp.MustCompile(`I{1,3}|\p{Lu}[\p{Ll} ]*-[\p{Lu}]?[\p{Ll}]*|\p{Lu}[\p{Ll}]*\.?`)

func FullName(s string) string {
	s = rxDropSpaces.ReplaceAllString(s, ``)
	return strings.Join(rxFindNames.FindAllString(s, -1), ` `)
}

// RegExp Get last name
var nameSuffix = []string{"esq", "esquire", "j", "jr", "jnr", "sr", "snr", "1st", "2nd", "3rd", "4th", "5th", "i", "ii", "iii", "iv", "v", "clu", "chfc", "cfp", "md", "phd", "jd", "llm", "do", "dc", "pc"}

func LastNameFirst(fullName string) (s string) {
	parts := strings.Split(fullName, " ")
	pLen := len(parts)
	switch {
	case pLen == 0:
		return ""
	case pLen < 2:
		return parts[0]
	default:
		hasSfx := false
		lastPart := strings.ToLower(strings.Replace(parts[pLen-1], ".", "", -1))
		for _, sfx := range nameSuffix {
			if sfx == lastPart {
				hasSfx = true
			}
		}
		if hasSfx {
			if pLen < 3 {
				return parts[pLen-2] + ", " + parts[pLen-1]
			}
			return parts[pLen-2] + ", " + parts[pLen-1] + strings.Join(parts[:pLen-2], " ")
		}
		return parts[pLen-1] + ", " + strings.Join(parts[:pLen-1], " ")
	}
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
