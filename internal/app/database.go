package app

import (
	"github.com/vinser/flibgolite/internal/core/config"
	"github.com/vinser/flibgolite/internal/store"
)

// InitDatabase initializes database connection and schema.
func (a *App) InitDatabase(cfg *config.Config) *store.DB {
	db := store.NewDB(cfg.Database.DSN)
	if !db.IsReady() {
		db.InitDB()
	}
	return db
}
