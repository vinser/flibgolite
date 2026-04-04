package interfaces

import (
	"github.com/vinser/flibgolite/internal/core/model"
)

// BookRepository - интерфейс для работы с книгами
type BookRepository interface {
	FindBookByID(id int64) (*model.Book, error)
	FindBooksByAuthor(authorID int64, limit, offset int) ([]*model.Book, error)
	FindBooksBySerie(serieID int64, limit, offset int) ([]*model.Book, error)
	SearchBooks(query string, limit, offset int) ([]*model.Book, int64, error)
	ListBooks(limit, offset int) ([]*model.Book, int64, error)
	GetBookCount() (int64, error)
}

// AuthorRepository - интерфейс для работы с авторами
type AuthorRepository interface {
	FindAuthorByID(id int64) (*model.Author, error)
	ListAuthors(limit, offset int) ([]*model.Author, int64, error)
	GetAuthorCount() (int64, error)
}

// SerieRepository - интерфейс для работы с сериями
type SerieRepository interface {
	FindSerieByID(id int64) (*model.Serie, error)
	ListSeries(limit, offset int) ([]*model.Serie, int64, error)
	GetSerieCount() (int64, error)
}

// GenreRepository - интерфейс для работы с жанрами
type GenreRepository interface {
	FindGenreByID(id int64) (*model.Genre, error)
	FindGenreByCode(code string) (*model.Genre, error)
	ListGenres(limit, offset int) ([]*model.Genre, int64, error)
	GetGenreCount() (int64, error)
}
