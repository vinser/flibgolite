package fb2

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/vinser/flibgolite/pkg/model"
	"github.com/vinser/flibgolite/pkg/parser"
)

func (fb *FB2) GetFormat() string {
	return "fb2"
}

func (fb *FB2) GetTitle() string {
	return strings.TrimSpace(fb.Description.TitleInfo.BookTitle)
}

func (fb *FB2) GetSort() string {
	return parser.GetSortTitle(fb.Description.TitleInfo.BookTitle, parser.GetLanguageTag(fb.Description.TitleInfo.Lang))
}

func (fb *FB2) GetYear() string {
	year := strconv.Itoa(fb.Description.PublishInfo.Year)
	if year == "" {
		year = fb.Description.TitleInfo.Date
	}
	rYear := []rune(year)
	if len(rYear) > 4 {
		rYear = rYear[len(rYear)-4:]
	}
	return strings.TrimSpace(string(rYear))
}

func (fb *FB2) GetPlot() string {
	return parser.StripHTMLTags(strings.Join(fb.Description.TitleInfo.Annotation.P, " "))
}

func (fb *FB2) GetCover() string {
	return strings.TrimPrefix(fb.Description.TitleInfo.CoverPage.Image.Href, "#")
}

func (fb *FB2) GetLanguage() *model.Language {
	return parser.GetLanguage(fb.Description.TitleInfo.Lang)
}

func (fb *FB2) GetAuthors() []*model.Author {
	authors := make([]*model.Author, 0, len(fb.Description.TitleInfo.Authors))
	if len(fb.Description.TitleInfo.Authors) == 1 &&
		fb.Description.TitleInfo.Authors[0].FirstName == "" &&
		fb.Description.TitleInfo.Authors[0].MiddleName == "" &&
		fb.Description.TitleInfo.Authors[0].LastName != "" &&
		strings.Contains(fb.Description.TitleInfo.Authors[0].LastName, ",") { // many authors are in the last name
		aLN := strings.Split(fb.Description.TitleInfo.Authors[0].LastName, ",")
		for _, a := range aLN {
			author := parser.AuthorByFullName(a)
			if author.Sort != "" {
				authors = append(authors, author)
			}
		}
		return authors
	}
	for _, a := range fb.Description.TitleInfo.Authors {
		author := parser.AuthorByFullName(fmt.Sprintf("%s %s %s", a.FirstName, a.MiddleName, a.LastName))
		if author.Sort != "" {
			authors = append(authors, author)
		}
	}
	if len(authors) == 0 {
		authors = append(authors,
			&model.Author{
				Name: "[author not specified]",
				Sort: "[author not specified]",
			},
		)
	}
	return authors
}

func (fb *FB2) GetGenres() []string {
	return fb.Description.TitleInfo.Genres
}

func (fb *FB2) GetKeywords() string {
	return fb.Description.TitleInfo.Keywords
}

func (fb *FB2) GetSerie() *model.Serie {
	if len(fb.Description.PublishInfo.Series) > 0 {
		return &model.Serie{Name: parser.Title(fb.Description.PublishInfo.Series[0].Name, fb.Description.TitleInfo.Lang)}
	} else if len(fb.Description.TitleInfo.Series) > 0 {
		return &model.Serie{Name: parser.Title(fb.Description.TitleInfo.Series[0].Name, fb.Description.TitleInfo.Lang)}
	} else {
		return &model.Serie{}
	}
}
func (fb *FB2) GetSerieNumber() int {
	if len(fb.Description.PublishInfo.Series) > 0 {
		return fb.Description.PublishInfo.Series[0].Number
	} else if len(fb.Description.TitleInfo.Series) > 0 {
		return fb.Description.TitleInfo.Series[0].Number
	} else {
		return 0
	}
}
