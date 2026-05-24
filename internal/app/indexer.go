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

// InitIndexerOnce initializes the indexer for one-time scanning.
func (a *App) InitIndexerOnce(cfg *config.Config, db *store.DB, genresTree *genres.GenresTree, stockLog *rlog.Log) *index.Handler {
	bookQueue := make(chan model.Book, cfg.Database.BOOK_QUEUE_SIZE)
	fileQueue := make(chan index.File, cfg.Database.FILE_QUEUE_SIZE)
	archive := make(map[string]*index.ArchiveTracker)

	stockHandler := &index.Handler{
		CFG:       cfg,
		LOG:       stockLog,
		DB:        db,
		GT:        genresTree,
		BookQueue: bookQueue,
		FileQueue: fileQueue,
		Archive:   archive,
	}

	stockHandler.StopDB = make(chan struct{})
	stockHandler.StopScan = make(chan struct{})
	stockHandler.Hashes = hash.InitHashes(db.DB)

	stockHandler.InitStockFolders()

	go stockHandler.ProcessBookQueue()
	for i := 0; i < cfg.Database.MAX_SCAN_THREADS; i++ {
		go stockHandler.ProcessFileQueue()
	}

	dir := cfg.Library.STOCK_DIR
	stockHandler.ScanDir(dir)

	return stockHandler
}

// InitIndexer initializes the indexer with background scanning.
func (a *App) InitIndexer(cfg *config.Config, db *store.DB, genresTree *genres.GenresTree, stockLog *rlog.Log) *index.Handler {
	bookQueue := make(chan model.Book, cfg.Database.BOOK_QUEUE_SIZE)
	fileQueue := make(chan index.File, cfg.Database.FILE_QUEUE_SIZE)
	archive := make(map[string]*index.ArchiveTracker)

	stockHandler := &index.Handler{
		CFG:       cfg,
		LOG:       stockLog,
		DB:        db,
		GT:        genresTree,
		BookQueue: bookQueue,
		FileQueue: fileQueue,
		Archive:   archive,
	}

	stockHandler.StopDB = make(chan struct{})
	stockHandler.StopScan = make(chan struct{})
	stockHandler.Hashes = hash.InitHashes(db.DB)

	stockHandler.InitStockFolders()

	go stockHandler.ProcessBookQueue()
	for i := 0; i < cfg.Database.MAX_SCAN_THREADS; i++ {
		go stockHandler.ProcessFileQueue()
	}

	go func() {
		defer func() { stockHandler.StopScan <- struct{}{} }()
		dir := cfg.Library.STOCK_DIR
		for {
			stockHandler.ScanDir(dir)
			time.Sleep(time.Duration(cfg.Database.POLL_DELAY) * time.Second)
			select {
			case <-stockHandler.StopScan:
				return
			default:
				continue
			}
		}
	}()

	return stockHandler
}
