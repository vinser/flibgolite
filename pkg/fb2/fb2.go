package fb2

import (
	"fmt"
	"strings"

	"github.com/vinser/flibgolite/pkg/model"
	"github.com/vinser/flibgolite/pkg/parser"
)

func (fb *FB2) GetFormat() string {
	return "fb2"
}

func (fb *FB2) GetTitle() string {
	return strings.TrimSpace(fb.Title)
}

func (fb *FB2) GetSort() string {
	return strings.ToUpper(strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(strings.TrimSpace(fb.Title), "An "), "A "), "The "))
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
	return strings.TrimSpace(string(rYear))
}

func (fb *FB2) GetPlot() string {
	return fb.Annotation.Text
}

func (fb *FB2) GetCover() string {
	return strings.TrimPrefix(fb.CoverPage.Href, "#")
}

func (fb *FB2) GetLanguage() *model.Language {
	base, _ := parser.GetLanguageTag(fb.Lang).Base()
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
		f := parser.RefineName(a.FirstName, fb.Lang)
		m := parser.RefineName(a.MiddleName, fb.Lang)
		l := parser.RefineName(a.LastName, fb.Lang)
		author.Name = parser.CollapceSpaces(fmt.Sprintf("%s %s %s", f, m, l))
		author.Sort = parser.CollapceSpaces(fmt.Sprintf("%s, %s %s", l, f, m))
		authors = append(authors, author)
	}
	return authors
}

func (fb *FB2) GetGenres() []string {
	return fb.Gengres
}

func (fb *FB2) GetKeywords() string {
	return fb.Keywords
}

func (fb *FB2) GetSerie() *model.Serie {
	return &model.Serie{Name: parser.Title(fb.Serie.Name, fb.Lang)}
}

func (fb *FB2) GetSerieNumber() int {
	return fb.Serie.Number
}
