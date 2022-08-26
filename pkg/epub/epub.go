package epub

import (
	"strings"

	"github.com/vinser/flibgolite/pkg/model"
	"github.com/vinser/flibgolite/pkg/parser"
	"golang.org/x/text/cases"
	"golang.org/x/text/language/display"
)

func (ep *OPF) GetFormat() string {
	return "epub"
}

func (ep *OPF) GetTitle() string {
	return strings.TrimSpace(strings.Join(ep.Metadata.Title, ", "))
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
		if cr.Role == "aut" || len(ep.Metadata.Creator) == 1 {
			fullName := parser.FullName(cr.Text)
			a.Name = fullName
			if cr.FileAs != "" {
				a.Sort = strings.Replace(parser.FullName(cr.FileAs), " ", ", ", 1)
			} else {
				a.Sort = parser.LastNameFirst(fullName)
			}
			authors = append(authors, a)
		}
	}
	return authors
}

func (ep *OPF) GetGenres() []string {
	return make([]string, 0)
}

func (ep *OPF) GetKeywords() string {
	return ""
}

func (ep *OPF) GetSerie() *model.Serie {
	return &model.Serie{}
}

func (ep *OPF) GetSerieNumber() int {
	return 0
}
