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
	FileOpenFailed
	FileIsNotRegular
)
const MIN_TITLEPLOT_LEN = 128

type BookHashes struct {
	Archives  map[string]map[string]int
	Files     map[string]int
	CRC32     map[uint32]int
	TitlePlot map[uint32]int
	mx        sync.RWMutex
}

func InitHashes(db *sqlx.DB) *BookHashes {
	count := 0
	db.QueryRowx(`SELECT count(*) FROM books`).Scan(&count)
	bh := &BookHashes{
		Archives:  make(map[string]map[string]int),
		Files:     make(map[string]int),
		CRC32:     make(map[uint32]int, count),
		TitlePlot: make(map[uint32]int, count),
	}
	rows, err := db.Query(`SELECT file, archive, crc32, title, plot FROM books`)
	if err != nil {
		log.Panicln(err)
	}
	defer rows.Close()
	for rows.Next() {
		b := &model.Book{}
		err := rows.Scan(&b.File, &b.Archive, &b.CRC32, &b.Title, &b.Plot)
		if err != nil {
			log.Panicln(err)
		}
		bh.Add(b.File, b.Archive)
		if b.CRC32 != 0 {
			bh.CRC32[b.CRC32] = 1
		}
		tpBytes := []byte(b.Title + b.Plot)
		if len(tpBytes) >= MIN_TITLEPLOT_LEN {
			bh.CRC32[crc32.ChecksumIEEE(tpBytes)] = 1
		}
	}
	return bh
}

func (bh *BookHashes) Add(file, archive string) {
	bh.mx.Lock()
	defer bh.mx.Unlock()
	if archive == "" {
		bh.Files[file] = 1
	} else {
		if _, ok := bh.Archives[archive]; !ok {
			bh.Archives[archive] = make(map[string]int)
		}
		if file != "" {
			bh.Archives[archive][file] = 1
		}
	}
}

func (bh *BookHashes) FileExists(file, archive string) bool {
	bh.mx.RLock()
	defer bh.mx.RUnlock()
	if archive == "" {
		_, ok := bh.Files[file]
		return ok
	}
	if _, ok := bh.Archives[archive]; ok {
		if _, ok := bh.Archives[archive][file]; ok {
			return true
		}
	}
	return false
}

func (bh *BookHashes) ArchiveExists(archive string) bool {
	bh.mx.RLock()
	defer bh.mx.RUnlock()
	_, ok := bh.Archives[archive]
	return ok
}

func (bh *BookHashes) GetState(b *model.Book, level string) BookState {
	bh.mx.Lock()
	defer bh.mx.Unlock()
	if b.Updated < 0 {
		return BookState(b.Updated)
	}
	if level != "N" {
		if _, ok := bh.CRC32[b.CRC32]; ok {
			return DuplicateCRC32
		}
		bh.CRC32[b.CRC32] = 1
	}
	if level == "S" {
		tpBytes := []byte(b.Title + b.Plot)
		if len(tpBytes) >= MIN_TITLEPLOT_LEN {
			tpCRC32 := crc32.ChecksumIEEE(tpBytes)
			if _, ok := bh.TitlePlot[tpCRC32]; ok {
				return DuplicateTitlePlot
			}
			bh.TitlePlot[tpCRC32] = 1
		}
	}
	return Unique
}
