package hash

import (
	"hash/crc32"
	"log"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/vinser/flibgolite/pkg/model"
)

type BookState int64

const (
	Unique BookState = -1 * iota
	DuplicateCRC32
	DuplicateTitlePlot
	FileIsEmpty
	FileHasErrors
	LanguageNotAccepted
	BadArchive
	UnsupportedFormat
)
const MIN_TITLEPLOT_LEN = 128

type BookHash struct {
}

type BookHashes struct {
	Archives  map[string]map[string]int
	Files     map[string]int
	CRC32     map[uint32]int
	TitlePlot map[uint32]int
	mx        sync.Mutex
}

func (bh *BookHashes) IsUnique(b *model.Book) BookState {
	bh.mx.Lock()
	defer bh.mx.Unlock()
	if _, ok := bh.CRC32[b.CRC32]; ok {
		return DuplicateCRC32
	}
	bh.CRC32[b.CRC32] = 1
	tpBytes := []byte(b.Title + b.Plot)
	if len(tpBytes) >= MIN_TITLEPLOT_LEN {
		tpCRC32 := crc32.ChecksumIEEE(tpBytes)
		if _, ok := bh.TitlePlot[tpCRC32]; ok {
			return DuplicateTitlePlot
		}
		bh.TitlePlot[tpCRC32] = 1
	}
	return Unique
}

// func (bh *BookHashes) SetState(crc uint32, title, plot string, bs BookState) {
// 	bh.mx.Lock()
// 	defer bh.mx.Unlock()
// 	switch bs {
// 	default:
// 		if crc != 0 {
// 			bh.CRC32[crc] = bs
// 		}
// 	case DuplicateTitlePlot:
// 		tpBytes := []byte(title + plot)
// 		if len(tpBytes) >= MIN_TITLEPLOT_LEN && crc != 0 {
// 			bh.CRC32[crc32.ChecksumIEEE(tpBytes)] = bs
// 		}
// 	}
// }

func InitHashes(db *sqlx.DB) *BookHashes {
	count := 0
	db.QueryRowx(`SELECT count(*) FROM books`).Scan(&count)
	bh := &BookHashes{
		Archives:  make(map[string]map[string]int),
		Files:     make(map[string]int),
		CRC32:     make(map[uint32]int, count),
		TitlePlot: make(map[uint32]int, count),
		mx:        sync.Mutex{},
	}
	rows, err := db.Query(`SELECT file, archive,crc32, title, plot FROM books`)
	if err != nil {
		log.Panicln(err)
	}
	defer rows.Close()
	bh.mx.Lock()
	defer bh.mx.Unlock()
	for rows.Next() {
		b := &model.Book{}
		err := rows.Scan(&b.File, &b.Archive, &b.CRC32, &b.Title, &b.Plot)
		if err != nil {
			log.Panicln(err)
		}
		if b.Archive == "" {
			bh.Files[b.File] = 1
		} else {
			if _, ok := bh.Archives[b.Archive]; !ok {
				bh.Archives[b.Archive] = make(map[string]int)
			}
			if b.File != "" {
				bh.Archives[b.Archive][b.File] = 1
			}
		}
		bh.Files[b.File+b.Archive] = 1
		bh.CRC32[b.CRC32] = 1
		tpBytes := []byte(b.Title + b.Plot)
		if len(tpBytes) >= MIN_TITLEPLOT_LEN {
			bh.CRC32[crc32.ChecksumIEEE(tpBytes)] = 1
		}
	}
	return bh
}
