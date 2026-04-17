package app

import (
	"github.com/vinser/flibgolite/internal/core/config"
	"github.com/vinser/flibgolite/internal/genres"
)

// InitGenres initializes the genres tree.
func (a *App) InitGenres(cfg *config.Config) *genres.GenresTree {
	return genres.NewGenresTree(cfg.Genres.TREE_FILE)
}
