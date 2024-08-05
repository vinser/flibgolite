package database

import (
	"database/sql"
	"fmt"

	"github.com/vinser/flibgolite/pkg/model"
)

func (db *DB) BookInfo(id int64) (*model.Book, error) {
	b := &model.Book{}
	q := `
		SELECT b.title, b.sort, b.plot, b.cover	FROM books as b	WHERE id=?`
	err := db.QueryRow(q, id).Scan(&b.Title, &b.Sort, &b.Plot, &b.Cover)
	if err != nil {
		if err == sql.ErrNoRows {
			return b, fmt.Errorf("book %d not found", id)
		}
		return b, fmt.Errorf("book %d not found: %w", id, err)
	}
	return b, nil
}

func (db *DB) BookLanguage(id int64) (*model.Language, error) {
	l := &model.Language{}
	q := `SELECT l.code, l.name FROM languages as l, books as b WHERE b.language_id=l.ID AND b.id=?`
	err := db.QueryRow(q, id).Scan(&l.Code, &l.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			l.Code = "en"
			return l, fmt.Errorf("book %d has no language set", id)
		}
		return l, fmt.Errorf("book %d has wrong language set: %w", id, err)
	}
	return l, nil
}

func (db *DB) BookAuthors(id int64) ([]*model.Author, error) {
	authors := []*model.Author{}
	q := `SELECT a.name, a.sort FROM authors as a, books_authors as ba WHERE ba.book_id=? AND ba.author_id=a.id`
	rows, err := db.Query(q, id)
	if err != nil {
		return authors, fmt.Errorf("book %d has authors error: %w", id, err)
	}
	defer rows.Close()

	for rows.Next() {
		a := &model.Author{}
		if err := rows.Scan(&a.Name, &a.Sort); err != nil {
			return authors, fmt.Errorf("book %d has authors error: %w", id, err)
		}
		authors = append(authors, a)
	}
	return authors, nil
}

func (db *DB) BookGenres(id int64) ([]string, error) {
	genres := []string{}
	q := `SELECT bg.genre_code FROM books_genres as bg WHERE bg.book_id=?`
	rows, err := db.Query(q, id)
	if err != nil {
		return genres, fmt.Errorf("book %d has genres error: %w", id, err)
	}
	defer rows.Close()

	for rows.Next() {
		g := ""
		if err := rows.Scan(&g); err != nil {
			return genres, fmt.Errorf("book %d has genres error: %w", id, err)
		}
		genres = append(genres, g)
	}
	return genres, nil
}
