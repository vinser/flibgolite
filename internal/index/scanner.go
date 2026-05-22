package index

import (
	"os"
	"path/filepath"
	"time"

	"github.com/vinser/flibgolite/internal/hash"
)

// ScanDir scans a directory for new books and processes them
func (h *Handler) ScanDir(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	entries, err := d.ReadDir(-1)
	if err != nil {
		return err
	}
	absDir, _ := filepath.Abs(dir)
	h.LOG.I.Printf("scanning folder %s for new books...\n", absDir)
	for _, entry := range entries {
		path, ext, err := h.isFileReady(dir, entry)
		if err != nil {
			h.LOG.I.Println(err)
			continue
		}
		switch {
		case ext == ".fb2":
			go func() {
				h.LOG.I.Println("file: ", entry.Name())
				err = h.parseFB2(path)
				if err != nil {
					h.LOG.W.Printf("Error processing file %s: %v", entry.Name(), err)
				}
			}()
		case ext == ".epub":
			go func() {
				h.LOG.I.Println("file: ", entry.Name())
				err = h.parseEPUB(path)
				if err != nil {
					h.LOG.W.Printf("Error processing file %s: %v", entry.Name(), err)
				}
			}()
		case ext == ".zip":
			start := time.Now()
			new := !h.Hashes.ArchiveExists(entry.Name())
			h.LOG.I.Println("zip: ", entry.Name())
			err = h.parseZipEntry(path)
			if err != nil {
				h.LOG.W.Println(err)
			}
			if new {
				h.LOG.S.Printf("%v elapsed for parsing %s ", time.Since(start), entry.Name())
			}
		default:
			h.LOG.D.Printf("file %s has not supported format \"%s\"\n", path, filepath.Ext(path))
			h.addFileToBookQueue(entry.Name(), "", hash.UnsupportedFormat)
		}
	}
	return nil
}
