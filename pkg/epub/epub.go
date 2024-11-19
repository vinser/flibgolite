package epub

import (
	"strings"
	"unicode"

	"github.com/vinser/flibgolite/pkg/model"
	"github.com/vinser/flibgolite/pkg/parser"
	"golang.org/x/text/cases"
	"golang.org/x/text/language/display"
)

func (ep *OPF) GetFormat() string {
	return "epub"
}

func (ep *OPF) GetTitle() string {
	if len(ep.Metadata.Title) > 0 {
		return strings.TrimSpace(ep.Metadata.Title[0])
	}
	return ""
}

func (ep *OPF) GetSort() string {
	l := ep.Lang
	if len(ep.Metadata.Language) > 0 {
		l = ep.Metadata.Language[0]
	}
	tag := parser.GetLanguageTag(l)
	title := ep.GetTitle()
	if base, _ := tag.Base(); base.String() == "en" {
		title = parser.Drop1stArticle(title)
	}
	return cases.Upper(tag).String(title)
}

func (ep *OPF) GetYear() string {
	return parser.PickYear(ep.Metadata.Date)
}

func (ep *OPF) GetPlot() string {
	return parser.StripHTMLTags(strings.Join(ep.Metadata.Description, " "))
}

func (ep *OPF) GetCover() string {
	if strings.TrimSpace(ep.Version) == "2.0" {
		content := ""
		for _, meta := range ep.Metadata.Meta {
			if meta.Name == "cover" {
				content = strings.TrimSpace(meta.Content)
				break
			}
		}
		if content == "" {
			return ""
		}
		for _, item := range ep.Manifest.Item {
			if item.ID == content {
				return strings.TrimSpace(item.Href)
			}
		}
		return ""
	}
	for _, item := range ep.Manifest.Item {
		if item.Properties == "cover-image" {
			return strings.TrimSpace(item.Href)
		}
	}
	return ""
}

func (ep *OPF) GetLanguage() *model.Language {
	l := ep.Lang
	if len(ep.Metadata.Language) > 0 {
		l = ep.Metadata.Language[0]
	}

	tag := parser.GetLanguageTag(l)
	base, _ := tag.Base()
	langName := cases.Title(tag).String(display.Self.Name(tag))

	return &model.Language{
		Code: base.String(),
		Name: langName,
	}
}

func (ep *OPF) GetAuthors() []*model.Author {
	authors := make([]*model.Author, 0)
	for _, cr := range ep.Metadata.Creator {
		a := &model.Author{}
		for _, meta := range ep.Metadata.Meta {
			if meta.Refines != "#"+cr.ID {
				continue
			}

			switch {
			case meta.Property == "role" && meta.Text == "aut":
				cr.Role = "aut"
			case meta.Property == "file-as":
				cr.FileAs = meta.Text
			}
		}
		if cr.Role == "aut" || cr.Role == "" || len(ep.Metadata.Creator) == 1 {
			parts := strings.Split(cr.Text, ",")
			name := parser.ParseFullName(parts[0])
			a.Name = strings.TrimSpace(strings.TrimSuffix(name.First+" "+name.Middle+" "+name.Last+" ("+name.Nick+")", " ()"))
			if cr.FileAs != "" {
				a.Sort = parser.AddCommaAfterLastName(parser.DelimitGluedName(cr.FileAs))
			} else {
				sortName := name.Last + ", " + name.First + " " + name.Middle + " (" + name.Nick + ")"
				a.Sort = strings.TrimSuffix(strings.TrimSpace(strings.TrimSuffix(sortName, " ()")), ",")
			}
			if len(a.Sort) > 0 {
				a.Sort = strings.ToUpper(a.Sort)
				authors = append(authors, a)
			}
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

func isSeparator(r rune) bool {
	return r == ',' || r == ';' || r == '-' || unicode.IsSpace(r)
}

func (ep *OPF) GetGenres() []string {
	return strings.FieldsFunc(strings.Join(ep.Metadata.Subject, " "), isSeparator)
}

func (ep *OPF) GetKeywords() string {
	return strings.Join(strings.FieldsFunc(strings.Join(ep.Metadata.Subject, " "), isSeparator), " ")
}

func (ep *OPF) GetSerie() *model.Serie {
	return &model.Serie{}
}

func (ep *OPF) GetSerieNumber() int {
	return 0
}
