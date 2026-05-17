package store

import (
	"database/sql"
	"fmt"
	"log"
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
