package index

import (
	"os"
	"path/filepath"
)

// moveFile moves a file to trash or stock directory based on error state
func (h *Handler) moveFile(filePath string, err error) {
	if err != nil && h.CFG.Library.TRASH_DIR != "" {
		os.Rename(filePath, filepath.Join(h.CFG.Library.TRASH_DIR, filepath.Base(filePath)))
		return
	}
	if filepath.Dir(filePath) == h.CFG.Library.STOCK_DIR {
		return
	}
	os.Rename(filePath, filepath.Join(h.CFG.Library.STOCK_DIR, filepath.Base(filePath)))
}
