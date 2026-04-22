package app

import (
	"github.com/vinser/flibgolite/internal/core/config"
)

// InitConfig loads application configuration from root directory.
func (a *App) InitConfig(rootDir string) *config.Config {
	return config.LoadConfig(rootDir)
}
