package app

import (
	"github.com/vinser/flibgolite/internal/core/config"
	"github.com/vinser/flibgolite/internal/rlog"
)

// InitLogs initializes application logs.
func (a *App) InitLogs(cfg *config.Config, withOpds bool) (stockLog, opdsLog *rlog.Log) {
	return cfg.InitLogs(withOpds)
}
