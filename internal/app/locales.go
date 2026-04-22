package app

import (
	"github.com/vinser/flibgolite/internal/core/config"
)

// InitLocales initializes locales from configuration file.
func (a *App) InitLocales(cfg *config.Config) {
	cfg.Locales.LoadLocales()
}
