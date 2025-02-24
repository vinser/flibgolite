package database

import (
	"database/sql"
	"log"

	"github.com/vinser/flibgolite/pkg/hash"
	"github.com/vinser/flibgolite/pkg/model"
)

func (tx *TX) PrepareStatements() {
	tx.Stmt["selectIdFromLanguages"] = tx.mustPrepare(`SELECT id FROM languages WHERE code=?`)
	tx.Stmt["insertIntoLanguages"] = tx.mustPrepare(`INSERT INTO languages (code, name) VALUES (?, ?)`)
	tx.Stmt["insertIntoBooks"] = tx.mustPrepare(`INSERT INTO books (file, crc32, archive, size, format, title, sort, year, language_id, plot, cover, keywords, serie_id, serie_num, updated) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	tx.Stmt["selectIdFromAuthors"] = tx.mustPrepare(`SELECT id FROM authors WHERE name=?`)
	tx.Stmt["insertIntoAuthors"] = tx.mustPrepare(`INSERT INTO authors (name, sort) VALUES (?, ?)`)
	tx.Stmt["insertIntoBooksAuthors"] = tx.mustPrepare(`INSERT INTO books_authors (book_id, author_id) VALUES (?, ?)`)
	tx.Stmt["insertIntoBooksGenres"] = tx.mustPrepare(`INSERT INTO books_genres (book_id, genre_code) VALUES (?, ?)`)
	tx.Stmt["selectIdFromSeries"] = tx.mustPrepare(`SELECT id FROM series WHERE name=?`)
	tx.Stmt["insertIntoSeries"] = tx.mustPrepare(`INSERT INTO series (name) VALUES (?)`)
}

// Books
func (tx *TX) NewBook(b *model.Book) {

	languageId := tx.NewLanguage(b.Language)
	serieId := tx.NewSerie(b.Serie)
	res, err := tx.Stmt["insertIntoBooks"].Exec(b.File, b.CRC32, b.Archive, b.Size, b.Format, b.Title, b.Sort, b.Year, languageId, b.Plot, b.Cover, b.Keywords, serieId, b.SerieNum, b.Updated)
	if err != nil {
		log.Panicln(err)
	}

	bookId, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		return
	}
	q := `INSERT INTO books_fts (rowid, title, keywords) VALUES (?, ?, ?)`
	_, err = tx.Exec(q, bookId, b.Title, b.Keywords)
	if err != nil {
		log.Panicln(err)
	}

	for _, author := range b.Authors {
		authorId := tx.NewAuthor(author)
		_, err = tx.Stmt["insertIntoBooksAuthors"].Exec(bookId, authorId)
	}
	if err != nil {
		log.Println(err)
	}

	for _, genre := range b.Genres {
		_, err = tx.Stmt["insertIntoBooksGenres"].Exec(bookId, genre)
	}
	if err != nil {
		log.Println(err)
	}
}

func (tx *TX) RecordBookState(b *model.Book, s hash.BookState) {
	_, err := tx.Stmt["insertIntoBooks"].Exec(b.File, b.CRC32, b.Archive, b.Size, b.Format, b.Title, b.Sort, b.Year, 0, b.Plot, b.Cover, b.Keywords, 0, b.SerieNum, int64(s))
	if err != nil {
		log.Panicln(err)
	}
}

// Languages
func (tx *TX) NewLanguage(l *model.Language) int64 {
	id := tx.FindLanguage(l)
	if id != 0 {
		return id
	}
	res, _ := tx.Stmt["insertIntoLanguages"].Exec(l.Code, l.Code)
	id, _ = res.LastInsertId()
	return id
}

func (tx *TX) FindLanguage(l *model.Language) int64 {
	var id int64 = 0
	err := tx.Stmt["selectIdFromLanguages"].QueryRow(l.Code).Scan(&id)
	if err == sql.ErrNoRows {
		return 0
	}
	return id
}

// Series
func (tx *TX) NewSerie(s *model.Serie) int64 {
	if s.Name == "" {
		return 0
	}
	id := tx.FindSerie(s)
	if id != 0 {
		return id
	}
	res, _ := tx.Stmt["insertIntoSeries"].Exec(s.Name)
	id, _ = res.LastInsertId()
	return id
}

func (tx *TX) FindSerie(s *model.Serie) int64 {
	var id int64 = 0
	err := tx.Stmt["selectIdFromSeries"].QueryRow(s.Name).Scan(&id)
	if err == sql.ErrNoRows {
		return 0
	}
	return id
}

// Authors
func (tx *TX) NewAuthor(a *model.Author) int64 {
	id := tx.FindAuthor(a)
	if id != 0 {
		return id
	}
	res, err := tx.Stmt["insertIntoAuthors"].Exec(a.Name, a.Sort)
	if err != nil {
		log.Printf("Name: %s, Sort: %s\n", a.Name, a.Sort)
		log.Panicln(err)
	}
	id, _ = res.LastInsertId()
	q := `INSERT INTO authors_fts (rowid, sort) VALUES (?, ?)`
	_, err = tx.Exec(q, id, a.Sort)
	if err != nil {
		log.Panicln(err)
	}
	return id
}

func (tx *TX) FindAuthor(a *model.Author) int64 {
	var id int64 = 0
	err := tx.Stmt["selectIdFromAuthors"].QueryRow(a.Name).Scan(&id)
	if err == sql.ErrNoRows {
		return 0
	}
	return id
}
