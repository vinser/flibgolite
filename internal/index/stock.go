package index

import (
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/vinser/flibgolite/internal/core/config"
	"github.com/vinser/flibgolite/internal/core/model"
	"github.com/vinser/flibgolite/internal/genres"
	"github.com/vinser/flibgolite/internal/hash"
	"github.com/vinser/flibgolite/internal/rlog"
	"github.com/vinser/flibgolite/internal/store"
)

type Handler struct {
	CFG       *config.Config
	Hashes    *hash.BookHashes
	DB        *store.DB
	GT        *genres.GenresTree
	LOG       *rlog.Log
	ScanWG    sync.WaitGroup
	FileQueue chan File
	BookQueue chan model.Book
	StopScan  chan struct{}
	StopDB    chan struct{}
}

type File struct {
	Reader  io.ReadCloser
	Name    string
	CRC32   uint32
	Archive string
	Size    int64
}

// InitStockFolders()
func (h *Handler) InitStockFolders() {
	if err := os.MkdirAll(h.CFG.Library.STOCK_DIR, 0776); err != nil {
		log.Fatalf("failed to create Library STOCK_DIR directory %s: %s", h.CFG.Library.STOCK_DIR, err)
	}
}

func (h *Handler) addFileToBookQueue(file, archive string, state hash.BookState) {
	h.BookQueue <- model.Book{
		File:    file,
		Archive: archive,
		Updated: int64(state),
	}
}

func (h *Handler) acceptLanguage(lang string) bool {
	if strings.Contains(h.CFG.ACCEPTED, "any") {
		return true
	}
	return strings.Contains(h.CFG.ACCEPTED, lang)
}
