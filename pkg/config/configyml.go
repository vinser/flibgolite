package config

const CONFIG_YML = `
library:
  # Selfexplained folders
  BOOK_STOCK: "books/stock"
  NEW_ACQUISITIONS: "books/new"
  TRASH: "books/trash"

language:
  # Locales folder. You can add your own locale file there like en.yml or ru.yml
  LOCALES: "config/locales"
  # Default english locale for opds feeds (bookreaders opds menu tree) can be changed to:
  # "ua" for Ukrainian, 
  # "ru" for Russian 
  DEFAULT: "en"

genres:
  TREE_FILE: "config/genres.xml"
  # Alternative genres tree can be used (Russian only, sorry) 
  # TREE_FILE: "config/alt_genres.xml"
  
database:
  DSN: "dbdata/flibgolite.db"
  # Delay before start each new acquisitions processing
  POLL_DELAY: 30 
  # Maximum simultaneous new aquisitios processing threads
  MAX_SCAN_THREADS: 3
  # Accept only these languages puplications. Add others if needed please.
  ACCEPTED_LANGS: "en,ua,ru"

logs:
  # Logs are here
  OPDS: "logs/opds.log"
  SCAN: "logs/scan.log"
  DEBUG: false

opds:
  # OPDS-server port so opds can be found at http://<server name or IP-address or localhost>:8085/opds
  PORT: 8085
  # OPDS feeds entries page size
  PAGE_SIZE: 30
`

var LOCALES_YML = map[string]string{
	`en`: `
alphabet: ABCDEFGHIJKLMNOPQRSTUVWXYZ
Book Authors: Book Authors
Choose an author of a book: Choose an author of a book
Book Genres: Book Genres
Choose a genre of a book: Choose a genre of a book
Book Series: Book Series
Choose a serie of a book: Choose a serie of a book
Choose from the found ones: Choose from the found ones
Titles: Titles
Authors: Authors
Found titles - %d: Found titles - %d
Total books - %d: Total books - %d
Found authors - %d: Found authors - %d
Alphabet: Alphabet
List books alphabetically: List books alphabetically
Series: Series
List books series: List books series
Genres: Genres
Book not found: Book not found
Total series - %d: Total series - %d
`,

	"ua": `
alphabet: АБВГҐДЕЄЖЗИІЇЙКЛМНОПРСТУФХЦЧШЩЬЮЯ
Book Authors: Автори
Choose an author of a book: Виберіть автора книги
Book Genres: Жанри
Choose a genre of a book: Виберіть жанр книги
Book Series: Серії
Choose a serie of a book: Виберіть серію книги
Choose from the found ones: Виберіть із знайдених
Titles: Книги
Authors: Автори
Found titles - %d: Знайдено книг - %d
Total books - %d: Книжок всього - %d
Found authors - %d: Знайдено авторів - %d
Alphabet: Алфавіт
List books alphabetically: Список книг по алфавіту
Series: Серії
List books series: Список книг за серіями
Genres: Жанри
Book not found: Книжка не знайдена
Total series - %d: Усього серій - %d
`,

	"ru": `
alphabet: АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ
Book Authors: Авторы
Choose an author of a book: Выбери автора книги
Book Genres: Жанры
Choose a genre of a book: Выбери жанр книги
Book Series: Серии
Choose a serie of a book: Выбери книжную серию
Choose from the found ones: Выбeри из найденных
Titles: Книги
Authors: Авторы
Found titles - %d: Найдено книг - %d
Total books - %d: Книг всего - %d
Found authors - %d: Найдено авторов - %d
Alphabet: Алфавит
List books alphabetically: Список книг по алфавиту
Series: Серии
List books series: Список книг по сериям
Genres: Жанры
Book not found: Книга не найдена
Total series - %d: Всего серий - %d
`,
}
