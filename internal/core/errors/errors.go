package errors

import "fmt"

var (
	// ErrBookNotFound - книга не найдена
	ErrBookNotFound = fmt.Errorf("book not found")
	// ErrAuthorNotFound - автор не найден
	ErrAuthorNotFound = fmt.Errorf("author not found")
	// ErrSerieNotFound - серия не найдена
	ErrSerieNotFound = fmt.Errorf("serie not found")
	// ErrGenreNotFound - жанр не найден
	ErrGenreNotFound = fmt.Errorf("genre not found")
	// ErrInvalidFormat - неверный формат файла
	ErrInvalidFormat = fmt.Errorf("invalid file format")
	// ErrParseError - ошибка парсинга
	ErrParseError = fmt.Errorf("parse error")
	// ErrConversionError - ошибка конвертации
	ErrConversionError = fmt.Errorf("conversion error")
	// ErrFileNotFound - файл не найден
	ErrFileNotFound = fmt.Errorf("file not found")
	// ErrPermissionDenied - отказ в доступе
	ErrPermissionDenied = fmt.Errorf("permission denied")
	// ErrDatabaseBusy - база данных занята
	ErrDatabaseBusy = fmt.Errorf("database busy")
)
