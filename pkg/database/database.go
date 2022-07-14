package database

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/vinser/flibgolite/pkg/model"

	// _ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)

type DB struct {
	*sql.DB
	mx sync.Mutex
}

// Books
func (db *DB) NewBook(b *model.Book) int64 {
	db.mx.Lock()
	defer db.mx.Unlock()

	bookId := db.FindBook(b)
	if bookId != 0 {
		return bookId
	}
	languageId := db.NewLanguage(b.Language)

	q := `
		INSERT INTO books (file, crc32, archive, size, format, title, sort, year,language_id, plot, cover, updated) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	res, err := db.Exec(q,
		b.File,
		b.CRC32,
		b.Archive,
		b.Size,
		b.Format,
		b.Title,
		b.Sort,
		b.Year,
		languageId,
		b.Plot,
		b.Cover,
		b.Updated,
	)
	if err != nil {
		log.Panic(err)
	}

	bookId, err = res.LastInsertId()
	if err != nil {
		log.Println(err)
		return 0
	}
	q = `INSERT INTO books_fts (rowid, title, keywords) VALUES (?, ?, ?)`
	_, err = db.Exec(q, bookId, b.Title, b.Keywords)
	if err != nil {
		log.Panicln(err)
	}

	for _, author := range b.Authors {
		authorId := db.NewAuthor(author)
		q = `INSERT INTO books_authors (book_id, author_id) VALUES (?, ?)`
		_, err = db.Exec(q, bookId, authorId)
	}
	if err != nil {
		log.Println(err)
	}

	for _, genre := range b.Genres {
		q = `INSERT INTO books_genres (book_id, genre_code) VALUES (?, ?)`
		_, err = db.Exec(q, bookId, genre)
	}
	if err != nil {
		log.Println(err)
	}

	serieId := db.NewSerie(b.Serie)
	if serieId != 0 {
		q = `INSERT INTO books_series (serie_num, book_id, serie_id) VALUES (?, ?, ?)`
		_, err = db.Exec(q, b.SerieNum, bookId, serieId)
		if err != nil {
			log.Println(err)
		}
	}

	return bookId
}

func (db *DB) FindBook(b *model.Book) int64 {
	var id int64 = 0
	q := `SELECT id FROM books WHERE sort LIKE ? and crc32=?`
	err := db.QueryRow(q, b.Sort, b.CRC32).Scan(&id)
	if err == sql.ErrNoRows {
		return 0
	}
	return id
}

func (db *DB) FindBookById(id int64) *model.Book {
	b := &model.Book{}
	q := `SELECT file, archive, format, title, cover FROM books WHERE id=?`
	err := db.QueryRow(q, id).Scan(&b.File, &b.Archive, &b.Format, &b.Title, &b.Cover)
	if err == sql.ErrNoRows {
		return nil
	}
	return b
}

func (db *DB) IsInStock(file string, crc32 uint32) bool {
	var id int64 = 0
	q := `SELECT id FROM books WHERE file=? AND crc32=?`
	err := db.QueryRow(q, file, crc32).Scan(&id)
	return err != sql.ErrNoRows
}

// Languages
func (db *DB) NewLanguage(l *model.Language) int64 {
	id := db.FindLanguage(l)
	if id != 0 {
		return id
	}
	q := `INSERT INTO languages (code, name) VALUES (?, ?)`
	res, _ := db.Exec(q, l.Code, l.Code)
	id, _ = res.LastInsertId()
	return id
}

func (db *DB) FindLanguage(l *model.Language) int64 {
	var id int64 = 0
	q := `SELECT id FROM languages WHERE code LIKE ?`
	err := db.QueryRow(q, l.Code).Scan(&id)
	if err == sql.ErrNoRows {
		return 0
	}
	return id
}

// Authors
func (db *DB) NewAuthor(a *model.Author) int64 {
	id := db.FindAuthor(a)
	if id != 0 {
		return id
	}
	q := `INSERT INTO authors (name, sort) VALUES (?, ?)`
	res, _ := db.Exec(q, a.Name, a.Sort)
	id, _ = res.LastInsertId()
	q = `INSERT INTO authors_fts (rowid, sort) VALUES (?, ?)`
	_, err := db.Exec(q, id, a.Sort)
	if err != nil {
		log.Panicln(err)
	}
	return id
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
			SELECT b.id, b.title, b.plot, b.cover, b.format 
			FROM books as b, books_authors as ba 
			WHERE ba.author_id=? AND b.id=ba.book_id ORDER BY b.sort
		`
		rows, err = db.pageQuery(q, limit, offset, authorId)
	} else {
		q = `
			SELECT b.id, b.title, b.plot, b.cover, b.format 
			FROM books as b, books_authors as ba, series as s, books_series as bs 
			WHERE ba.author_id=? AND ba.book_id=b.id AND bs.book_id=b.id AND bs.serie_id=? GROUP BY b.title
		`
		rows, err = db.pageQuery(q, limit, offset, authorId, serieId)
	}
	if err != nil {
		log.Println("DB page query error: ", err.Error())
	}
	defer rows.Close()
	books := []*model.Book{}

	for rows.Next() {
		b := &model.Book{}
		if err = rows.Scan(&b.ID, &b.Title, &b.Plot, &b.Cover, &b.Format); err != nil {
			log.Fatal(err)
		}
		books = append(books, b)
	}
	return books
}

func (db *DB) AuthorBookSeries(id int64) []*model.Serie {
	series := []*model.Serie{}
	q := `
		SELECT s.id, s.name 
		FROM books_authors as ba, books as b, books_series as bs, series as s 
		WHERE ba.author_id=? AND b.id=ba.book_id AND b.id=bs.book_id AND s.id=bs.serie_id GROUP BY s.name
	`
	rows, err := db.Query(q, id)
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

func (db *DB) AuthorByID(id int64) *model.Author {
	author := &model.Author{}
	q := `SELECT name, sort FROM authors WHERE id=?`
	err := db.QueryRow(q, id).Scan(&author.Name, &author.Sort)
	if err == sql.ErrNoRows {
		return nil
	}
	return author
}

func (db *DB) FindAuthor(a *model.Author) int64 {
	var id int64 = 0
	q := `SELECT id FROM authors WHERE sort LIKE ?`
	err := db.QueryRow(q, a.Sort).Scan(&id)
	if err == sql.ErrNoRows {
		return 0
	}
	return id
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
		SELECT b.id, b.title, b.plot, b.cover, b.format 
		FROM books as b, books_genres as bg 
		WHERE bg.genre_code=? AND b.id=bg.book_id ORDER BY b.sort
	`
	rows, err := db.pageQuery(q, limit, offset, genreCode)
	if err != nil {
		log.Println("DB page query error: ", err.Error())
	}
	defer rows.Close()
	books := []*model.Book{}

	for rows.Next() {
		b := &model.Book{}
		if err = rows.Scan(&b.ID, &b.Title, &b.Plot, &b.Cover, b.Format); err != nil {
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
func (db *DB) NewSerie(s *model.Serie) int64 {
	if s.Name == "" {
		return 0
	}
	id := db.FindSerie(s)
	if id != 0 {
		return id
	}
	q := `INSERT INTO series (name) VALUES (?)`
	res, _ := db.Exec(q, s.Name)
	id, _ = res.LastInsertId()
	return id
}

func (db *DB) ListSerieBooks(id int64, limit, offset int) []*model.Book {
	q := `
		SELECT b.id, b.title, b.plot, b.cover, b.format 
		FROM books as b, books_series as bs 
		WHERE bs.serie_id=? AND b.id=bs.book_id ORDER BY bs.serie_num
	`
	rows, err := db.pageQuery(q, limit, offset, id)
	if err != nil {
		log.Println("DB page query error: ", err.Error())
	}
	defer rows.Close()
	books := []*model.Book{}

	for rows.Next() {
		b := &model.Book{}
		if err = rows.Scan(&b.ID, &b.Title, &b.Plot, &b.Cover, &b.Format); err != nil {
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
			(SELECT serie_id, book_id, count(*) as c FROM books_series GROUP BY serie_id HAVING c>2) as bs,
			books as b,
			languages as l 
			WHERE 
			sr.id=bs.serie_id AND 
			s IN (`, abc, `) AND
			b.id=bs.book_id AND
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
			(SELECT serie_id, book_id, count(*) as c FROM books_series GROUP BY serie_id HAVING c>2) as bs,
			books as b,
			languages as l 
			WHERE 
			sr.id=bs.serie_id AND 
			sn LIKE ? AND
			b.id=bs.book_id AND
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
		FROM series as s, books_series as bs, books as b, languages as l 
		WHERE 
		s.name LIKE ? AND 
		s.id=bs.serie_id AND
		b.id=bs.book_id AND
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

func (db *DB) SerieByID(id int64) *model.Serie {
	serie := &model.Serie{}
	q := `SELECT name FROM series WHERE id=?`
	err := db.QueryRow(q, id).Scan(&serie.Name)
	if err == sql.ErrNoRows {
		return nil
	}
	return serie
}

func (db *DB) FindSerie(s *model.Serie) int64 {
	var id int64 = 0
	q := `SELECT id FROM series WHERE name LIKE ?`
	err := db.QueryRow(q, s.Name).Scan(&id)
	if err == sql.ErrNoRows {
		return 0
	}
	return id
}

// Search
func (db *DB) SearchBooksCount(pattern string) int64 {
	var c int64 = 0
	q := `SELECT count(*) as c FROM books_fts WHERE title MATCH ? OR keywords MATCH ?`
	err := db.QueryRow(q, pattern, pattern).Scan(&c)
	if err == sql.ErrNoRows {
		return 0
	}
	return c
}

func (db *DB) PageFoundBooks(pattern string, limit, offset int) []*model.Book {
	q := `
		WITH s AS(
			SELECT rowid FROM books_fts WHERE title MATCH ? OR keywords MATCH ? ORDER BY rank DESC LIMIT ? OFFSET ?
		)
		SELECT b.id, b.title, b.plot, b.cover, b.format FROM books AS b, s WHERE b.id=s.rowid
	`
	rows, err := db.Query(q, pattern, pattern, limit, offset)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	books := []*model.Book{}

	for rows.Next() {
		b := &model.Book{}
		if err := rows.Scan(&b.ID, &b.Title, &b.Plot, &b.Cover, &b.Format); err != nil {
			log.Fatal(err)
		}
		books = append(books, b)
	}
	return books
}

func (db *DB) SearchAuthorsCount(pattern string) int64 {
	var c int64 = 0
	q := `SELECT count(*) as c FROM authors_fts WHERE sort MATCH ?`
	err := db.QueryRow(q, "^"+pattern).Scan(&c)
	if err == sql.ErrNoRows {
		return 0
	}
	return c
}

func (db *DB) PageFoundAuthors(pattern string, limit, offset int) []*model.Author {
	q := `
		WITH s AS(
			SELECT rowid FROM authors_fts WHERE sort MATCH ? LIMIT ? OFFSET ?
		)
		SELECT a.id, a.name, a.sort, count(*) as c FROM authors AS a, books_authors AS ba, s 
		WHERE a.id=s.rowid AND a.id=ba.author_id GROUP BY a.sort ORDER BY c DESC
	`
	rows, err := db.Query(q, "^"+pattern, limit, offset)
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

// ==================================
func NewDB(dsn string) *DB {
	err := os.MkdirAll(filepath.Dir(dsn), 0775)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
	// db, err := sql.Open("sqlite3", dsn)
	db, err := sql.Open("sqlite", dsn+"?_pragma=busy_timeout%3d10000")
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(10)
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	return &DB{
		DB: db,
		mx: sync.Mutex{},
	}
}

func (db *DB) InitDB() {
	if !db.IsReady() {
		db.execFile(SQLITE_DB_INIT)
	}
}

func (db *DB) DropDB() {
	if db.IsReady() {
		db.execFile(SQLITE_DB_DROP)
	}
}

func (db *DB) IsReady() bool {
	var err error
	rows, err := db.Query(`SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%' AND name NOT LIKE 'test%'`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	return rows.Next()
}

func (db *DB) execFile(sql string) {
	scanner := bufio.NewScanner(strings.NewReader(sql))
	scanner.Split(bufio.ScanLines)
	q := ""

	for scanner.Scan() {
		q += scanner.Text()
		if strings.Contains(q, ";") {
			_, err := db.Exec(q)
			q = ""
			if err != nil {
				log.Fatal(err)
			}
		}
	}
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
