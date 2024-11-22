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
    keywords TEXT,
    serie_id INTEGER,
    serie_num INTEGER,
    updated INTEGER
);
-- CREATE UNIQUE INDEX book_crc32_idx ON books (crc32);  -- crc32 is unique?
CREATE INDEX book_crc32_idx ON books (crc32); 
CREATE INDEX book_file_idx ON books (file);
CREATE INDEX book_archive_idx ON books (archive);
CREATE INDEX book_title_idx ON books (title);
CREATE INDEX book_sort_idx ON books (sort COLLATE NOCASE);
CREATE INDEX book_language_idx ON books (language_id);
CREATE INDEX book_serie_idx ON books (serie_id);
CREATE INDEX book_updated_idx ON books (updated);

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
