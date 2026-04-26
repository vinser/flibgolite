package app

import (
	"github.com/vinser/flibgolite/internal/core/config"
	"github.com/vinser/flibgolite/internal/store"
)

// InitDatabase initializes database connection and schema.
func (a *App) InitDatabase(cfg *config.Config) (*store.DB, error) {
	db, err := store.NewDB(cfg.Database.DSN)
	if err != nil {
		return nil, err
	}
	ready, err := db.IsReady()
	if err != nil {
		return nil, err
	}
	if !ready {
		db.InitDB()
	}
	return db, nil
}
