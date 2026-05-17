package store

import (
	"database/sql"
	"log"
	"time"

	"github.com/vinser/flibgolite/internal/core/model"
)

// Latest

func (db *DB) LatestBooksCount(days int) int64 {
	var c int64 = 0
	q := `SELECT count(*) as c FROM books WHERE updated > ?`
	err := db.QueryRow(q, time.Now().Unix()-int64(days*24*60*60)).Scan(&c)
	if err == sql.ErrNoRows {
		return 0
	}
	return c
}

func (db *DB) PageLatestBooks(days, limit, offset int) []*model.Book {
	q := `
	SELECT b.id, b.file, b.archive, b.size, b.format, b.title, b.year, b.plot, b.cover, ifnull(s.name, ''), b.serie_num, ifnull(l.code, '') 
	FROM books as b
	LEFT JOIN series as s ON b.serie_id=s.id
	JOIN languages as l ON b.language_id=l.id
	WHERE b.updated > ? 
	ORDER BY b.id DESC
	`
	rows, err := db.pageQuery(q, limit, offset, time.Now().Unix()-int64(days*24*60*60))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	books := []*model.Book{}
	for rows.Next() {
		b := &model.Book{
			Language: &model.Language{},
			Serie:    &model.Serie{},
		}
		if err := rows.Scan(&b.ID, &b.File, &b.Archive, &b.Size, &b.Format, &b.Title, &b.Year, &b.Plot, &b.Cover, &b.Serie.Name, &b.SerieNum, &b.Language.Code); err != nil {
			log.Fatal(err)
		}
		books = append(books, b)
	}
	return books
}
