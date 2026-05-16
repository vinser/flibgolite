package store

import (
	"database/sql"
	"fmt"
	"log"
	"unicode/utf8"

	"github.com/vinser/flibgolite/internal/core/model"
)

// Authors

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
		SELECT id, name, sort, 
			(
				SELECT COUNT(*) 
				FROM books_authors 
				WHERE author_id = authors.id
			) AS count
		FROM authors
		WHERE sort LIKE ?
		ORDER BY sort
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
