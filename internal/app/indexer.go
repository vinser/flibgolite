package app

import (
	"time"

	"github.com/vinser/flibgolite/internal/core/config"
	"github.com/vinser/flibgolite/internal/core/model"
	"github.com/vinser/flibgolite/internal/genres"
	"github.com/vinser/flibgolite/internal/hash"
	"github.com/vinser/flibgolite/internal/index"
	"github.com/vinser/flibgolite/internal/rlog"
	"github.com/vinser/flibgolite/internal/store"
)

// initIndexerCommon initializes the shared indexer components.
func (a *App) initIndexerCommon(cfg *config.Config, db *store.DB, genresTree *genres.GenresTree, stockLog *rlog.Log) *index.Handler {
	bookQueue := make(chan model.Book, cfg.Database.BOOK_QUEUE_SIZE)
	fileQueue := make(chan index.File, cfg.Database.FILE_QUEUE_SIZE)
	archive := make(map[string]*index.ArchiveTracker)

	h := &index.Handler{
		CFG:       cfg,
		LOG:       stockLog,
		DB:        db,
		GT:        genresTree,
		BookQueue: bookQueue,
		FileQueue: fileQueue,
		Archive:   archive,
	}

	h.StopDB = make(chan struct{})
	h.StopScan = make(chan struct{})
	h.Hashes = hash.InitHashes(db.DB)

	h.InitStockFolders()

	// Start workers
	go h.ProcessBookQueue()
	for i := 0; i < cfg.Database.MAX_SCAN_THREADS; i++ {
		go h.ProcessFileQueue()
	}

	return h
}

// InitIndexerOnce initializes the indexer for one-time scanning.
func (a *App) InitIndexerOnce(cfg *config.Config, db *store.DB, genresTree *genres.GenresTree, stockLog *rlog.Log) *index.Handler {
	h := a.initIndexerCommon(cfg, db, genresTree, stockLog)

	// One-time scan
	dir := cfg.Library.STOCK_DIR
	h.ScanDir(dir)

	return h
}

// InitIndexer initializes the indexer with background scanning.
func (a *App) InitIndexer(cfg *config.Config, db *store.DB, genresTree *genres.GenresTree, stockLog *rlog.Log) *index.Handler {
	h := a.initIndexerCommon(cfg, db, genresTree, stockLog)

	// Background scanning loop
	go func() {
		defer func() { h.StopScan <- struct{}{} }()
		dir := cfg.Library.STOCK_DIR

		for {
			h.ScanDir(dir)
			time.Sleep(time.Duration(cfg.Database.POLL_DELAY) * time.Second)

			select {
			case <-h.StopScan:
				return
			default:
				continue
			}
		}
	}()

	return h
}
