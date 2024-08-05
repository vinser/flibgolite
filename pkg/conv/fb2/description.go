package fb2

import (
	"github.com/vinser/flibgolite/pkg/conv/epub2"
)

func (p *FB2Parser) parseDescription(e *epub2.EPUB) error {
	p.Skip()
	book, err := p.DB.BookInfo(p.BookId)
	if err != nil {
		p.LOG.E.Println(err)
		return err
	}
	book.Language, err = p.DB.BookLanguage(p.BookId)
	if err != nil {
		p.LOG.E.Println(err)
	}
	book.Authors, err = p.DB.BookAuthors(p.BookId)
	if err != nil {
		p.LOG.E.Println(err)
	}
	book.Genres, err = p.DB.BookGenres(p.BookId)
	if err != nil {
		p.LOG.E.Println(err)
	}
	e.AddMetadataLanguage(book.Language.Code)
	e.AddMetadataTitle(book.Title)
	if book.Plot != "" {
		e.AddMetadataDescription(book.Plot)
	}
	if book.Cover != "" {
		e.AddMetadataCover(book.Cover)
		if err := e.AddItem("cover", "cover", `<div class="cover"><img class="coverimage" alt="Cover" src="`+book.Cover+`" /></div>`); err != nil {
			return err
		}
	}
	for _, a := range book.Authors {
		e.AddMetadataAuthor(a.Name, a.Sort)
	}
	for _, g := range book.Genres {
		e.AddMetadataSubject(g)
	}

	return nil
}
