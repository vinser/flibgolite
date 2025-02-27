package database

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/vinser/flibgolite/pkg/model"
	"github.com/vinser/flibgolite/pkg/user"
	"golang.org/x/crypto/bcrypt"
)

// Auth
var ErrorBadUserOrPassword = fmt.Errorf("Bad user name or password")

func (db *DB) GetUserByUsername(username string) (user.User, error) {
	// TODO Replace with real request
	users := map[string]user.User{
		"admin": {
			ID:       0,
			Username: "admin",
			Password: "admin",
			Email:    "admin@localhost",
		},
		"john": {
			ID:       1,
			Username: "john",
			Password: "p@ss",
			Email:    "john@localhost",
		},
	}

	if u, ok := users[username]; ok {
		hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
		if err != nil {
			return user.User{}, err
		}
		u.Password = string(hashed)
		return u, nil
	} else {
		return user.User{}, ErrorBadUserOrPassword
	}
}

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
	prefixLen := utf8.RuneCountInString(prefix) + 1
	var (
		rows *sql.Rows
		err  error
		q    string
	)
	if prefixLen == 1 {
		if abc != "" {
			q = `
			SELECT id, name, SUBSTR(sort,1,1) as s, COUNT(*) 
			FROM authors 
			WHERE s IN(` + abc + `) 
			GROUP BY s
		`
		} else {
			q = `
			SELECT id, name, SUBSTR(sort,1,1) as s, COUNT(*) 
			FROM authors 
			WHERE sort NOT LIKE '[author not specified]' 
			GROUP BY s
			`
		}
		rows, err = db.Query(q)
	} else {
		q = fmt.Sprint(`
			SELECT id, name, SUBSTR(sort,1,`, fmt.Sprint(prefixLen), `) as s, COUNT(*)
			FROM authors 
			WHERE sort LIKE ? 
			GROUP BY s
			`)
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
	if len(authors) == 1 && authors[0].Count > 1 {
		pref := string([]rune(authors[0].Sort)[:prefixLen])
		return db.ListAuthors(pref, abc)
	}
	return authors
}

func (db *DB) AuthorNotSpecifiedId() int64 {
	q := `SELECT id FROM authors WHERE sort LIKE '[author not specified]'`
	var id int64
	err := db.QueryRow(q).Scan(&id)
	if err != nil {
		return 0
	}
	return id
}

func (db *DB) ListAuthorWithTotals(prefix string) []*model.Author {
	authors := []*model.Author{}
	q := `
		SELECT a.id, a.name, a.sort, count(*) 
		FROM authors as a, books_authors as ba 
		WHERE sort LIKE ? AND a.id=ba.author_id 
		GROUP BY a.sort
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

// Authors

func (db *DB) AuthorBookSeries(authorId int64) []*model.Serie {
	series := []*model.Serie{}
	q := `
		SELECT b.serie_id, s.name
		FROM books_authors as ba 
		JOIN books as b ON b.id=ba.book_id
		JOIN series as s ON s.id=b.serie_id
		WHERE ba.author_id=? AND b.serie_id!=0
		GROUP BY b.serie_id
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
		SELECT a.id, a.name, a.sort 
		FROM authors as a 
		JOIN books_authors as ba ON a.id=ba.author_id 
		WHERE ba.book_id=?
		ORDER BY a.sort
	`
	rows, err := db.Query(q, bookId)
	if err != nil {
		return authors
	}
	defer rows.Close()

	for rows.Next() {
		a := &model.Author{}
		if err := rows.Scan(&a.ID, &a.Name, &a.Sort); err != nil {
			return authors
		}
		authors = append(authors, a)
	}
	return authors
}

// Genres

func (db *DB) PageGenreBooks(genreCode string, limit, offset int) []*model.Book {
	q := `
		SELECT b.id, b.file, b.archive, b.size, b.format, b.title, b.sort, b.year, b.plot, b.cover,  ifnull(s.name, ''), b.serie_num, ifnull(l.code, '') 
		FROM books as b 
		LEFT JOIN books_genres as bg ON b.id=bg.book_id
		LEFT JOIN series as s ON b.serie_id=s.id
		LEFT JOIN languages as l ON b.language_id=l.id
		WHERE bg.genre_code=? 
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
		SELECT 
			SUBSTR(s.name, 1, 1) AS p,
			COUNT(DISTINCT s.id)
		FROM series AS s
		JOIN books AS b ON s.id = b.serie_id
		JOIN languages AS l ON l.id = b.language_id 
		WHERE l.code LIKE ? AND p IN(` + abc + `)
		GROUP BY p
		`)
		rows, err = db.Query(q, lang+"%", prefix+"%")
	} else {
		q := `
		SELECT 
			SUBSTR(s.name, 1, ?) AS p,
			COUNT(DISTINCT s.id)
		FROM series AS s
		JOIN books AS b ON s.id = b.serie_id
		JOIN languages AS l ON l.id = b.language_id 
		WHERE l.code LIKE ? AND s.name LIKE ?
		GROUP BY p
		`
		rows, err = db.Query(q, prefixLen, lang+"%", prefix+"%")
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
	SELECT 
		s.id, s.name, COUNT(b.id)
	FROM series AS s
	JOIN books AS b ON s.id = b.serie_id
	JOIN languages AS l ON l.id = b.language_id 
	WHERE l.code LIKE ? AND s.name LIKE ?
	GROUP BY s.id
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
