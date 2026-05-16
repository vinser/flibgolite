package store

import (
	"database/sql"
	"log"

	"github.com/vinser/flibgolite/internal/core/model"
)

// Books

func (db *DB) FindBookById(id int64) *model.Book {
	b := &model.Book{}
	q := `SELECT file, archive, format, title, cover FROM books WHERE id=?`
	err := db.QueryRow(q, id).Scan(&b.File, &b.Archive, &b.Format, &b.Title, &b.Cover)
	if err == sql.ErrNoRows {
		return nil
	}
	return b
}

func (db *DB) CountLanguageBooks(languageCode string) int64 {
	var c int64 = 0
	q := `SELECT count(*) FROM languages as l, books as b WHERE l.code=? AND b.language_id=l.ID`
	err := db.QueryRow(q, languageCode).Scan(&c)
	if err == sql.ErrNoRows {
		return 0
	}
	return c
}

func (db *DB) ListAuthorBooks(authorId, serieId int64, limit, offset int) []*model.Book {
	var (
		q    string
		rows *sql.Rows
		err  error
	)
	if serieId == 0 {
		q = `
		SELECT b.id, b.file, b.archive, b.size, b.format, b.title, b.year, b.plot, b.cover, ifnull(s.name, ''), b.serie_num, ifnull(l.code, '') 
		FROM books as b 
		JOIN books_authors as ba ON b.id=ba.book_id 
		LEFT JOIN series as s ON b.serie_id=s.id
		JOIN languages as l ON b.language_id=l.id
		WHERE ba.author_id=?  
		ORDER BY b.sort
		`
		rows, err = db.pageQuery(q, limit, offset, authorId)
	} else {
		q = `
		SELECT b.id, b.file, b.archive, b.size, b.format, b.title, b.year, b.plot, b.cover, s.name, b.serie_num, ifnull(l.code, '') 
		FROM books as b
		JOIN books_authors as ba ON b.id=ba.book_id
		JOIN series as s ON b.serie_id=s.id
		JOIN languages as l ON b.language_id=l.id
		WHERE ba.author_id=? AND b.serie_id=?
		GROUP BY b.title 
		ORDER BY b.serie_num
		`
		rows, err = db.pageQuery(q, limit, offset, authorId, serieId)
	}
	if err != nil {
		log.Println("DB page query error: ", err.Error())
	}
	defer rows.Close()
	books := []*model.Book{}

	for rows.Next() {
		b := &model.Book{
			Language: &model.Language{},
			Serie:    &model.Serie{},
		}
		if err = rows.Scan(&b.ID, &b.File, &b.Archive, &b.Size, &b.Format, &b.Title, &b.Year, &b.Plot, &b.Cover, &b.Serie.Name, &b.SerieNum, &b.Language.Code); err != nil {
			log.Fatal(err)
		}
		books = append(books, b)
	}
	return books
}
