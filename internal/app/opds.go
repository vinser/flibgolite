package app

import (
	"fmt"
	"net/http"

	"github.com/vinser/flibgolite/internal/core/config"
	"github.com/vinser/flibgolite/internal/genres"
	"github.com/vinser/flibgolite/internal/opds"
	"github.com/vinser/flibgolite/internal/rlog"
	"github.com/vinser/flibgolite/internal/store"
	"golang.org/x/text/message"
)

// InitOPDS initializes OPDS handler and HTTP server (without starting it).
func (a *App) InitOPDS(cfg *config.Config, db *store.DB, genresTree *genres.GenresTree, opdsLog *rlog.Log) (*opds.Handler, *http.Server) {
	opdsHandler := &opds.Handler{
		CFG: cfg,
		LOG: opdsLog,
		DB:  db,
		GT:  genresTree,
		MP:  make(map[string]*message.Printer, len(cfg.Locales.Languages)),
	}

	for k, v := range cfg.Locales.Languages {
		opdsHandler.MP[k] = message.NewPrinter(v.Tag)
	}

	auth := opdsHandler.NewAuth()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.OPDS.PORT),
		Handler: auth(opdsHandler),
	}

	return opdsHandler, server
}
