package store

import (
	"database/sql"
	"log"

	"github.com/vinser/flibgolite/internal/core/model"
)

// Genres

func (db *DB) PageGenreBooks(genreCode string, limit, offset int) []*model.Book {
	q := `
		SELECT b.id, b.file, b.archive, b.size, b.format, b.title, b.sort, b.year, b.plot, b.cover,  ifnull(s.name, ''), b.serie_num, ifnull(l.code, '') 
		FROM books_genres AS bg
		JOIN books AS b ON bg.book_id = b.id
		LEFT JOIN series AS s ON b.serie_id = s.id
		LEFT JOIN languages AS l ON b.language_id = l.id
		WHERE bg.genre_code = ? 
		ORDER BY b.sort
		`
	rows, err := db.pageQuery(q, limit, offset, genreCode)
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
		if err = rows.Scan(&b.ID, &b.File, &b.Archive, &b.Size, &b.Format, &b.Title, &b.Sort, &b.Year, &b.Plot, &b.Cover, &b.Serie.Name, &b.SerieNum, &b.Language.Code); err != nil {
			log.Fatal(err)
		}
		books = append(books, b)
	}
	return books
}

func (db *DB) CountGenreBooks(genreCode string) int64 {
	var c int64 = 0
	q := `SELECT count(*) FROM books_genres as bg WHERE bg.genre_code=?`
	err := db.QueryRow(q, genreCode).Scan(&c)
	if err == sql.ErrNoRows {
		return 0
	}
	return c
}
