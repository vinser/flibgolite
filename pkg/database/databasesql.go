package database

const SQLITE_DB_INIT = `
DROP TABLE IF EXISTS archives;
CREATE TABLE archives (
    id INTEGER PRIMARY KEY,
    name TEXT,
    commited INTEGER
);
CREATE UNIQUE INDEX archives_name_idx ON archives (name);

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
    archive_id INTEGER,
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
CREATE UNIQUE INDEX book_crc_idx ON books (crc32);
CREATE INDEX book_file_idx ON books (file);
CREATE INDEX book_archive_idx ON books (archive_id);
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
const SQLITE_DB_DROP = `
DROP TABLE IF EXISTS languages;
DROP TABLE IF EXISTS authors;
DROP TABLE IF EXISTS books;
DROP TABLE IF EXISTS genres;
DROP TABLE IF EXISTS series;
DROP TABLE IF EXISTS books_authors;
DROP TABLE IF EXISTS books_genres;
DROP TABLE IF EXISTS books_series;
`
