package app

import (
	"github.com/vinser/flibgolite/internal/core/config"
)

// App represents the main application structure
type App struct{}

// New creates a new App instance
func New(cfg *config.Config) *App {
	return &App{}
}
