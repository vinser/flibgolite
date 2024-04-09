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
	Count int // for intermediate keeping author book counts
}

type Archive struct {
	ID       int64
	Name     string
	Commited int64
}

type Book struct {
	ID       int64
	File     string
	CRC32    uint32
	Archive  *Archive
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
	Keywords string
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
	Count int // for intermediate keeping serie book counts
}
