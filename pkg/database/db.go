package database

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/vinser/flibgolite/pkg/config"
	"github.com/vinser/flibgolite/pkg/model"
	"github.com/vinser/flibgolite/pkg/rlog"

	_ "modernc.org/sqlite"
)

const (
	SQLITE_DB_BUSY_TIMEOUT = 10000
	SQLITE_DB_INIT         = `
DROP TABLE IF EXISTS languages;
CREATE TABLE languages (
    id INTEGER PRIMARY KEY,
    code TEXT,
    name TEXT
);
CREATE UNIQUE INDEX languages_code_idx ON languages (code);
CREATE INDEX languages_name_idx ON languages (name);

DROP TABLE IF EXISTS authors;
CREATE TABLE authors (
    id INTEGER PRIMARY KEY,
    name TEXT,
    sort TEXT
);
CREATE UNIQUE INDEX authots_name_idx ON authors (name);
CREATE INDEX authots_sort_idx ON authors (sort COLLATE NOCASE);

DROP TABLE IF EXISTS authors_fts;
CREATE VIRTUAL TABLE authors_fts USING fts5(sort, content='', tokenize='unicode61 remove_diacritics 2');

DROP TABLE IF EXISTS books;
CREATE TABLE books (
    id INTEGER PRIMARY KEY,
    file TEXT,
    crc32 INTEGER,
    archive TEXT,
    size INTEGER,
    format TEXT,
    title TEXT,
    sort TEXT,
    year TEXT,
    language_id INTEGER,
    plot TEXT,
    cover TEXT,
    updated INTEGER
);
CREATE UNIQUE INDEX book_crc32_idx ON books (crc32);
CREATE INDEX book_file_idx ON books (file);
CREATE INDEX book_archive_idx ON books (archive);
CREATE INDEX book_title_idx ON books (title);
CREATE INDEX book_sort_idx ON books (sort COLLATE NOCASE);
CREATE INDEX book_language_idx ON books (language_id);

DROP TABLE IF EXISTS books_fts;
CREATE VIRTUAL TABLE books_fts USING fts5(title, keywords, content='', tokenize='unicode61 remove_diacritics 2');

DROP TABLE IF EXISTS series;
CREATE TABLE series (
    id INTEGER PRIMARY KEY,
    name TEXT
);
CREATE UNIQUE INDEX series_name_idx ON series (name);

DROP TABLE IF EXISTS books_authors;
CREATE TABLE books_authors (
    id INTEGER PRIMARY KEY,
    book_id INTEGER,
    author_id INTEGER
);
CREATE INDEX books_authors_book_idx ON books_authors (book_id);
CREATE INDEX books_authors_author_idx ON books_authors (author_id);

DROP TABLE IF EXISTS books_genres;
CREATE TABLE books_genres (
    id INTEGER PRIMARY KEY,
    book_id INTEGER,
    genre_code TEXT
);
CREATE INDEX books_genres_genre_code_idx ON books_genres (genre_code);
CREATE INDEX books_genres_book_idx ON books_genres (book_id);

DROP TABLE IF EXISTS books_series;
CREATE TABLE books_series (
    id INTEGER PRIMARY KEY,
    serie_num INTEGER DEFAULT 0,
    book_id INTEGER,
    serie_id INTEGER
);
CREATE INDEX books_series_book_idx ON books_series (book_id);
CREATE INDEX books_series_serie_idx ON books_series (serie_id);
`
	SQLITE_DB_DROP = `
DROP TABLE IF EXISTS languages;
DROP TABLE IF EXISTS authors;
DROP TABLE IF EXISTS books;
DROP TABLE IF EXISTS genres;
DROP TABLE IF EXISTS series;
DROP TABLE IF EXISTS books_authors;
DROP TABLE IF EXISTS books_genres;
DROP TABLE IF EXISTS books_series;
`
)

// ==================================
type Handler struct {
	CFG   *config.Config
	DB    *DB
	TX    *TX
	LOG   *rlog.Log
	WG    *sync.WaitGroup
	Queue <-chan model.Book
	Stop  chan struct{}
}

type DB struct {
	*sqlx.DB
}

// ==================================
func NewDB(dsn string) *DB {
	err := os.MkdirAll(filepath.Dir(dsn), 0775)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
	// db, err := sql.Open("sqlite", dsn+"?_pragma=busy_timeout(10000)&_pragma=journal_mode(wal)")
	options := fmt.Sprintf("?_pragma=busy_timeout(%d)&_pragma=journal_mode(delete)", SQLITE_DB_BUSY_TIMEOUT)
	db, err := sqlx.Open("sqlite", dsn+options)

	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(30)
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	DB := &DB{
		DB: db,
	}

	return DB
}

func (db *DB) Close() {
	db.DB.Close()
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

// ==================================
type TX struct {
	*sqlx.Tx
	Stmt map[string]*sqlx.Stmt
}

func (db *DB) txBegin() *TX {
	TX := &TX{
		Tx:   db.DB.MustBegin(),
		Stmt: map[string]*sqlx.Stmt{},
	}
	TX.PrepareStatements()
	return TX
}

func (tx *TX) txEnd() {
	err := tx.Tx.Commit()
	if err != nil && !errors.Is(err, sql.ErrTxDone) {
		log.Printf("Commit failed: %v", err)
		if err = tx.Tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			log.Printf("Rollback failed: %v", err)
		}
	}
	for _, stmt := range tx.Stmt {
		stmt.Close()
	}
}

func (tx *TX) mustPrepare(query string) *sqlx.Stmt {
	stmt, err := tx.Tx.Preparex(query)
	if err != nil {
		panic(err)
	}
	return stmt
}
