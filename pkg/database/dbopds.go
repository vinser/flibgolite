package database

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/vinser/flibgolite/pkg/model"
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

func (db *DB) ListAuthors(prefix, abc string) []*model.Author {
	l := utf8.RuneCountInString(prefix) + 1
	var (
		rows *sql.Rows
		err  error
	)
	if l == 1 {
		q := fmt.Sprint(`SELECT id, name, substr(sort,1,1) as s, count(*) as c FROM authors WHERE s IN(`, abc, `) GROUP BY s`)

		rows, err = db.Query(q)
	} else {
		q := fmt.Sprint(`SELECT id, name, substr(sort,1,`, fmt.Sprint(l), `) as s, count(*) as c FROM authors WHERE sort LIKE ? GROUP BY s`)
		rows, err = db.Query(q, prefix+"%")
	}
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	authors := []*model.Author{}

	for rows.Next() {
		var a *model.Author = &model.Author{}
		if err := rows.Scan(&a.ID, &a.Name, &a.Sort, &a.Count); err != nil {
			log.Fatal(err)
		}
		authors = append(authors, a)
	}
	return authors
}

func (db *DB) ListAuthorWithTotals(prefix string) []*model.Author {
	authors := []*model.Author{}
	q := `
		SELECT a.id, a.name, a.sort, count(*) 
		FROM authors as a, books_authors as ba 
		WHERE sort LIKE ? AND a.id=ba.author_id GROUP BY a.sort
	`
	rows, err := db.Query(q, prefix+"%")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var a *model.Author = &model.Author{}
		if err := rows.Scan(&a.ID, &a.Name, &a.Sort, &a.Count); err != nil {
			log.Fatal(err)
		}
		authors = append(authors, a)
	}
	return authors
}

func (db *DB) ListAuthorBooks(authorId, serieId int64, limit, offset int) []*model.Book {
	var (
		q    string
		rows *sql.Rows
		err  error
	)
	if serieId == 0 {
		q = `
		SELECT b.id, b.size, b.format, b.title, b.year, b.plot, b.cover,  ifnull(s.name, ''), b.serie_num 
		FROM books as b, books_authors as ba
		LEFT JOIN series as s ON b.serie_id=s.id
		WHERE ba.author_id=? AND b.id=ba.book_id ORDER BY b.sort
		`
		rows, err = db.pageQuery(q, limit, offset, authorId)
	} else {
		q = `
		SELECT b.id, b.size, b.format, b.title, b.year, b.plot, b.cover, s.name, b.serie_num 
		FROM books as b, books_authors as ba, series as s
		WHERE ba.author_id=? AND ba.book_id=b.id AND b.serie_id=? AND b.serie_id=s.id GROUP BY b.title ORDER BY b.serie_num
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
		if err = rows.Scan(&b.ID, &b.Size, &b.Format, &b.Title, &b.Year, &b.Plot, &b.Cover, &b.Serie.Name, &b.SerieNum); err != nil {
			log.Fatal(err)
		}
		books = append(books, b)
	}
	return books
}

// Authors

func (db *DB) AuthorBookSeries(authorId int64) []*model.Serie {
	series := []*model.Serie{}
	q := `
		SELECT b.serie_id, s.name
		FROM books_authors as ba, books as b, series as s
		WHERE ba.author_id=? AND b.id=ba.book_id AND b.serie_id!=0 AND b.serie_id=s.id GROUP BY b.serie_id
	`
	rows, err := db.Query(q, authorId)
	if err != nil {
		return series
	}
	defer rows.Close()

	for rows.Next() {
		s := &model.Serie{}
		if err := rows.Scan(&s.ID, &s.Name); err != nil {
			return series
		}
		series = append(series, s)
	}
	return series
}

func (db *DB) AuthorByID(authorId int64) *model.Author {
	author := &model.Author{}
	q := `SELECT name, sort FROM authors WHERE id=?`
	err := db.QueryRow(q, authorId).Scan(&author.Name, &author.Sort)
	if err == sql.ErrNoRows {
		return nil
	}
	return author
}

func (db *DB) AuthorsByBookId(bookId int64) []*model.Author {
	authors := []*model.Author{}
	q := `
		SELECT a.id, a.name 
		FROM authors as a, books_authors as ba 
		WHERE ba.book_id=? AND ba.author_id=a.id ORDER BY a.sort
	`
	rows, err := db.Query(q, bookId)
	if err != nil {
		return authors
	}
	defer rows.Close()

	for rows.Next() {
		a := &model.Author{}
		if err := rows.Scan(&a.ID, &a.Name); err != nil {
			return authors
		}
		authors = append(authors, a)
	}
	return authors
}

// Genres

func (db *DB) ListGenreBooks(genreCode string, limit, offset int) []*model.Book {
	q := `
		SELECT b.id, b.size, b.format, b.title, b.sort, b.year, b.plot, b.cover,  ifnull(s.name, ''), b.serie_num 
		FROM books as b, books_genres as bg
		LEFT JOIN series as s ON b.serie_id=s.id
		WHERE bg.genre_code=? AND b.id=bg.book_id ORDER BY b.sort
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
		if err = rows.Scan(&b.ID, &b.Size, &b.Format, &b.Title, &b.Sort, &b.Year, &b.Plot, &b.Cover, &b.Serie.Name, &b.SerieNum); err != nil {
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

// Series

func (db *DB) ListSerieBooks(id int64, limit, offset int) []*model.Book {
	q := `
		SELECT b.id, b.size, b.format, b.title, b.year, b.plot, b.cover, s.name, b.serie_num 
		FROM books as b, series as s
		WHERE b.serie_id=? AND b.serie_id=s.id ORDER BY b.serie_num
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
		if err = rows.Scan(&b.ID, &b.Size, &b.Format, &b.Title, &b.Year, &b.Plot, &b.Cover, &b.Serie.Name, &b.SerieNum); err != nil {
			log.Fatal(err)
		}
		books = append(books, b)
	}
	return books
}

// select id, substr(name,1,1) as s, count(*) as c FROM series group by s order by name<'а', `name`<'a',`name`;

func (db *DB) ListSeries(prefix, lang, abc string) []*model.Serie {
	l := utf8.RuneCountInString(prefix) + 1
	var (
		rows *sql.Rows
		err  error
	)
	if l == 1 {
		q := fmt.Sprint(`
			SELECT sr.id, substr(sr.name,1,1) as s, count(*) as c 
			FROM 
			series as sr, 
			(SELECT serie_id, language_id, count(*) as c FROM books GROUP BY serie_id HAVING c>2) as b,
			languages as l 
			WHERE 
			sr.id=b.serie_id AND 
			s IN (`, abc, `) AND
			l.id=b.language_id AND
			l.code=?
			GROUP BY s
			`)
		rows, err = db.Query(q, lang)
	} else {
		q := fmt.Sprint(`
			SELECT sr.id, substr(sr.name,1,`, fmt.Sprint(l), `) as sn, count(*) as c 
			FROM 
			series as sr, 
			(SELECT serie_id, language_id, count(*) as c FROM books GROUP BY serie_id HAVING c>2) as b,
			languages as l 
			WHERE 
			sr.id=b.serie_id AND 
			sn LIKE ? AND
			l.id=b.language_id AND
			l.code=?
			GROUP BY sn
					`)
		rows, err = db.Query(q, prefix+"%", lang)
	}
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	series := []*model.Serie{}
	for rows.Next() {
		a := &model.Serie{}
		if err := rows.Scan(&a.ID, &a.Name, &a.Count); err != nil {
			log.Fatal(err)
		}
		series = append(series, a)
	}
	return series
}

func (db *DB) ListSeriesWithTotals(prefix, lang string) []*model.Serie {
	series := []*model.Serie{}
	q := `
		SELECT s.id, s.name, count(*) as c 
		FROM series as s, books as b, languages as l 
		WHERE 
		s.name LIKE ? AND 
		s.id=b.serie_id AND
		l.id=b.language_id AND
		l.code=?
		GROUP BY s.name HAVING c>2
	`
	rows, err := db.Query(q, prefix+"%", lang)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		s := &model.Serie{}
		if err := rows.Scan(&s.ID, &s.Name, &s.Count); err != nil {
			log.Fatal(err)
		}
		series = append(series, s)
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

// Search
func (db *DB) SearchBooksCount(pattern string) int64 {
	var c int64 = 0
	q := `SELECT count(*) as c FROM books_fts WHERE title MATCH ?`
	err := db.QueryRow(q, pattern).Scan(&c)
	if err == sql.ErrNoRows {
		return 0
	}
	return c
}

func (db *DB) PageFoundBooks(pattern string, limit, offset int) []*model.Book {
	foundIDs := func(pattern string, limit, offset int) []string {
		q := `SELECT rowid FROM books_fts WHERE title MATCH ? ORDER BY rank DESC LIMIT ? OFFSET ?`
		rows, err := db.Query(q, pattern, limit, offset)
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
	}(pattern, limit, offset)
	q := `
	SELECT b.id, b.size, b.format, b.title, b.year, b.plot, b.cover, ifnull(s.name, ''), b.serie_num 
	FROM books as b
	LEFT JOIN series as s ON b.serie_id=s.id
	WHERE b.id IN ( ` + strings.Join(foundIDs, ",") + ` ) 
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
		if err := rows.Scan(&b.ID, &b.Size, &b.Format, &b.Title, &b.Year, &b.Plot, &b.Cover, &b.Serie.Name, &b.SerieNum); err != nil {
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
	// err := db.QueryRow(q, "^"+pattern).Scan(&c)
	err := db.QueryRow(q, pattern).Scan(&c)
	if err == sql.ErrNoRows {
		return 0
	}
	return c
}

func (db *DB) PageFoundAuthors(pattern string, limit, offset int) []*model.Author {
	q := `
	SELECT a.id, a.name, a.sort, count(*) as c FROM authors AS a, books_authors AS ba 
	WHERE a.id=ba.author_id AND a.id in (SELECT rowid FROM authors_fts WHERE sort MATCH ?)
	GROUP BY a.sort ORDER BY a.sort LIMIT ? OFFSET ?
	`
	// rows, err := db.Query(q, "^"+pattern, limit, offset)
	rows, err := db.Query(q, pattern, limit, offset)
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
