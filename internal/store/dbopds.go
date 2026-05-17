package store

import (
	"database/sql"
	"log"
	"strconv"
	"strings"
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

// Search

const (
	SearchBookByTitleMode   = "title"
	SearchBookByKeywordMode = "keywords"
)

func (db *DB) SearchBooksCountByTitle(pattern string) int64 {
	return db.searchBooksCount(SearchBookByTitleMode, pattern)
}

func (db *DB) SearchBooksCountByKeyword(pattern string) int64 {
	return db.searchBooksCount(SearchBookByKeywordMode, pattern)
}

func (db *DB) searchBooksCount(mode, pattern string) int64 {
	var c int64 = 0
	q := `SELECT count(*) as c FROM books_fts WHERE ` + mode + ` MATCH ?`
	err := db.QueryRow(q, pattern).Scan(&c)
	if err == sql.ErrNoRows {
		return 0
	}
	return c
}

func (db *DB) PageFoundBooksByTitle(pattern string, limit, offset int) []*model.Book {
	return db.pageFoundBooks(SearchBookByTitleMode, pattern, limit, offset)
}

func (db *DB) PageFoundBooksByKeywords(pattern string, limit, offset int) []*model.Book {
	return db.pageFoundBooks(SearchBookByKeywordMode, pattern, limit, offset)
}

func (db *DB) pageFoundBooks(mode, pattern string, limit, offset int) []*model.Book {
	foundIDs := func(mode, pattern string, limit, offset int) []string {
		q := `SELECT rowid 
			FROM books_fts 
			WHERE ` + mode + ` MATCH ? 
			ORDER BY rank 
			`
		rows, err := db.pageQuery(q, limit, offset, pattern)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		foundIDs := []string{}
		for rows.Next() {
			var id int
			if err := rows.Scan(&id); err != nil {
				log.Fatal(err)
			}
			foundIDs = append(foundIDs, strconv.Itoa(id))
		}
		return foundIDs
	}(mode, pattern, limit, offset)
	q := `
	SELECT b.id, b.file, b.archive, b.size, b.format, b.title, b.year, b.plot, b.cover, ifnull(s.name, ''), b.serie_num, ifnull(l.code, '') 
	FROM books as b
	LEFT JOIN series as s ON b.serie_id=s.id
	LEFT JOIN languages as l ON b.language_id=l.id
	WHERE b.id IN (` + strings.Join(foundIDs, ",") + `) 
	`
	rows, err := db.Query(q)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	booksIdx := map[int64]*model.Book{}
	for rows.Next() {
		b := &model.Book{
			Language: &model.Language{},
			Serie:    &model.Serie{},
		}
		if err := rows.Scan(&b.ID, &b.File, &b.Archive, &b.Size, &b.Format, &b.Title, &b.Year, &b.Plot, &b.Cover, &b.Serie.Name, &b.SerieNum, &b.Language.Code); err != nil {
			log.Fatal(err)
		}
		booksIdx[b.ID] = b
	}

	books := []*model.Book{}
	for _, b := range foundIDs {
		i, _ := strconv.ParseInt(b, 10, 64)
		books = append(books, booksIdx[i])
	}
	return books
}

func (db *DB) SearchAuthorsCount(pattern string) int64 {
	var c int64 = 0
	q := `SELECT count(*) as c FROM authors_fts WHERE sort MATCH ?`
	// err := db.QueryRow(q, pattern).Scan(&c)
	err := db.QueryRow(q, "^"+pattern).Scan(&c)
	if err == sql.ErrNoRows {
		return 0
	}
	return c
}

func (db *DB) PageFoundAuthors(pattern string, limit, offset int) []*model.Author {
	q := `
	SELECT a.id, a.name, a.sort, count(*) as c 
	FROM authors AS a
	JOIN books_authors AS ba ON a.id=ba.author_id 
	WHERE a.id in (SELECT rowid FROM authors_fts WHERE sort MATCH ?)
	GROUP BY a.sort 
	ORDER BY a.sort 
	`
	// rows, err := db.pageQuery(q, limit, offset, pattern)
	rows, err := db.pageQuery(q, limit, offset, "^"+pattern)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	authors := []*model.Author{}

	for rows.Next() {
		a := &model.Author{}
		if err := rows.Scan(&a.ID, &a.Name, &a.Sort, &a.Count); err != nil {
			log.Fatal(err)
		}
		authors = append(authors, a)
	}
	return authors
}

func (db *DB) pageQuery(query string, limit, offset int, args ...interface{}) (*sql.Rows, error) {
	if limit > 0 {
		query += " LIMIT ?"
		args = append(args, limit)
	}
	if offset > 0 {
		query += " OFFSET ?"
		args = append(args, offset)
	}
	// log.Println("query: ", query, " args: ", args)
	rows, err := db.Query(query, args...)
	return rows, err
}
