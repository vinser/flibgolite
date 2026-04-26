package index

import (
	"time"

	"github.com/vinser/flibgolite/internal/core/model"
	"github.com/vinser/flibgolite/internal/hash"
	"github.com/vinser/flibgolite/internal/parsers"
)

// createBookFromParser creates a model.Book from parser data
func (h *Handler) createBookFromParser(p parsers.Parser, file string, archive string, size int64, crc32 uint32) *model.Book {
	return &model.Book{
		File:     file,
		CRC32:    crc32,
		Archive:  archive,
		Size:     size,
		Format:   p.GetFormat(),
		Title:    p.GetTitle(),
		Sort:     p.GetSort(),
		Year:     p.GetYear(),
		Plot:     p.GetPlot(),
		Cover:    p.GetCover(),
		Language: p.GetLanguage(),
		Authors:  p.GetAuthors(),
		Genres:   p.GetGenres(),
		Keywords: p.GetKeywords(),
		Serie:    p.GetSerie(),
		SerieNum: p.GetSerieNumber(),
		Updated:  time.Now().UnixNano(),
	}
}

// processLanguage checks if the book language is accepted and returns error if not
func (h *Handler) processLanguage(p parsers.Parser, file, archive string) error {
	language := p.GetLanguage()
	if !h.acceptLanguage(language.Code) {
		h.addFileToBookQueue(file, archive, hash.LanguageNotAccepted)
		return &LanguageNotAcceptedError{Language: language.Code, File: file}
	}
	return nil
}

// LanguageNotAcceptedError represents an error when book language is not accepted
type LanguageNotAcceptedError struct {
	Language string
	File     string
}

func (e *LanguageNotAcceptedError) Error() string {
	return "publication language \"" + e.Language + "\" is configured as not accepted, file " + e.File + " has been skipped"
}
