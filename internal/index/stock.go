package index

import (
	"io"
	"os"
	"strings"
	"sync"
	"time"

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
	Archive   map[string]*ArchiveTracker
}

type File struct {
	Reader  io.ReadCloser
	Name    string
	CRC32   uint32
	Archive string
	Size    int64
}

type ArchiveTracker struct {
	Total   int64 // total FB2 entries in ZIP
	Pending int64 // how many books still not written to DB
}

// InitStockFolders()
func (h *Handler) InitStockFolders() {
	if err := os.MkdirAll(h.CFG.Library.STOCK_DIR, 0776); err != nil {
		h.LOG.E.Printf("failed to create Library STOCK_DIR directory %s: %s", h.CFG.Library.STOCK_DIR, err)
		os.Exit(1)
	}
}

func (h *Handler) WaitForQueuesToDrain() {
	// Wait until FileQueue and BookQueue become empty
	for len(h.FileQueue) > 0 || len(h.BookQueue) > 0 {
		time.Sleep(50 * time.Millisecond)
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
