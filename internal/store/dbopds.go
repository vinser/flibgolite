package store

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/vinser/flibgolite/internal/core/model"
)

// Series

func (db *DB) ListSerieBooks(id int64, limit, offset int) []*model.Book {
	q := `
		SELECT b.id, b.file, b.archive, b.size, b.format, b.title, b.year, b.plot, b.cover, s.name, b.serie_num, ifnull(l.code, '') 
		FROM books as b 
		JOIN series as s ON s.id=b.serie_id
		JOIN languages as l ON b.language_id=l.id
		WHERE b.serie_id=?
		ORDER BY b.serie_num
	`
	rows, err := db.pageQuery(q, limit, offset, id)
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

func (db *DB) ListSeries(prefix, lang, abc string) []*model.Serie {
	prefixLen := utf8.RuneCountInString(prefix) + 1
	var (
		rows *sql.Rows
		err  error
	)
	if prefixLen == 1 && abc != "" {
		q := fmt.Sprint(`
		SELECT SUBSTR(name, 1, 1) AS p,	COUNT(id)
		FROM series
		WHERE SUBSTR(name, 1, 1) IN(` + abc + `)
		  AND EXISTS (
			  SELECT 1 FROM books AS b
			  JOIN languages AS l ON l.id = b.language_id 
			  WHERE b.serie_id = series.id AND l.code LIKE ?
		  )
		GROUP BY p
		`)
		rows, err = db.Query(q, lang+"%")
	} else {
		q := `
		SELECT SUBSTR(name, 1, ?) AS p,	COUNT(id)
		FROM series
		WHERE name LIKE ?
		  AND EXISTS (
			  SELECT 1 FROM books AS b
			  JOIN languages AS l ON l.id = b.language_id 
			  WHERE b.serie_id = series.id AND l.code LIKE ?
		  )
		GROUP BY p
		`
		rows, err = db.Query(q, prefixLen, prefix+"%", lang+"%")
	}

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	series := []*model.Serie{}
	for rows.Next() {
		a := &model.Serie{}
		if err := rows.Scan(&a.Name, &a.Count); err != nil {
			log.Fatal(err)
		}
		series = append(series, a)
	}
	if len(series) == 1 && series[0].Count > 1 {
		pref := string([]rune(series[0].Name)[:prefixLen])
		return db.ListSeries(pref, lang, abc)
	}
	return series
}

func (db *DB) ListSeriesWithTotals(prefix, lang string) []*model.Serie {
	series := []*model.Serie{}
	q := `
	SELECT s.id, s.name, 
		(
			SELECT COUNT(b.id) 
			FROM books AS b 
			JOIN languages AS l ON l.id = b.language_id 
			WHERE b.serie_id = s.id AND l.code LIKE ?
		) AS count
	FROM series AS s
	WHERE s.name LIKE ?
	ORDER BY s.name
	`
	rows, err := db.Query(q, lang+"%", prefix+"%")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		s := &model.Serie{}
		if err := rows.Scan(&s.ID, &s.Name, &s.Count); err != nil {
			log.Fatal(err)
		}
		if s.Count > 0 {
			series = append(series, s)
		}
	}
	return series
}

func (db *DB) SerieByID(serieId int64) *model.Serie {
	serie := &model.Serie{}
	q := `SELECT name FROM series WHERE id=?`
	err := db.QueryRow(q, serieId).Scan(&serie.Name)
	if err == sql.ErrNoRows {
		return nil
	}
	return serie
}

func (db *DB) SerieByBookID(bookId int64) *model.Serie {
	serie := &model.Serie{}
	q := `SELECT s.id, s.name FROM books as b JOIN series as s ON s.id=b.serie_id WHERE b.id=?`
	err := db.QueryRow(q, bookId).Scan(&serie.ID, &serie.Name)
	if err == sql.ErrNoRows {
		return nil
	}
	return serie
}

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
