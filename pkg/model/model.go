package model

type Language struct {
	ID   int64
	Code string
	Name string
}

type Author struct {
	ID    int64
	Name  string
	Sort  string
	Count int // for intermediate storing author book counts
}

type Book struct {
	ID       int64
	File     string
	CRC32    uint32
	Archive  string
	Size     int64
	Format   string
	Title    string
	Sort     string
	Year     string
	Plot     string
	Cover    string
	Language *Language
	Authors  []*Author
	Genres   []string
	Serie    *Serie
	SerieNum int
	Updated  int64
}

type Genre struct {
	ID    int64
	Code  string
	Bunch string
	Name  string
}

type Serie struct {
	ID    int64
	Name  string
	Count int
}
