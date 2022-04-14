package parser

import "github.com/vinser/flibgolite/pkg/model"

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
	GetSerie() *model.Serie
	GetSerieNumber() int
}
