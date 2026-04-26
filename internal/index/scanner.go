package index

import (
	"archive/zip"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/vinser/flibgolite/internal/hash"
)

// isFileReady checks if a file is ready for processing
func (h *Handler) isFileReady(dir string, ent fs.DirEntry) (path string, ext string, err error) {
	info, err := ent.Info()
	if err != nil {
		return "", "", err
	}
	if info.Mode().IsRegular() {
		path = filepath.Join(dir, info.Name())
		ext = strings.ToLower(filepath.Ext(info.Name()))
		oldSize := info.Size()
		poll := time.Microsecond * 100
		wait := time.Second * 10
		for {
			time.Sleep(poll)
			info, err = ent.Info()
			if err != nil {
				return "", "", err
			}
			if info.Size() == oldSize {
				if info.Size() == 0 {
					err := fmt.Errorf("file %s has size of zero", path)
					h.addFileToBookQueue(info.Name(), "", hash.FileIsEmpty)
					h.moveFile(path, err)
					return "", "", err
				}
				// check if file is ready
				h.LOG.D.Println("Check if file is not busy", path)
				switch ext {
				case ".zip", ".epub":
					for {
						time.Sleep(poll)
						r, err := zip.OpenReader(path)
						if err == nil {
							r.Close()
							h.LOG.D.Println("Final polling period for file", path, ":", poll)
							return path, ext, nil
						}
						poll *= 2
						if poll > wait {
							return "", "", fmt.Errorf("file %s is busy and postponed until the next scan", path)
						}
					}
				default:
					time.Sleep(poll)
					return path, ext, nil
				}
			}
			oldSize = info.Size()
		}
	}
	h.addFileToBookQueue(info.Name(), "", hash.FileIsNotRegular)
	return "", "", fmt.Errorf("file %s is not a regular file", path)
}

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
				h.moveFile(path, err)
				if err != nil {
					h.LOG.W.Println(err)
				}
			}()
		case ext == ".epub":
			go func() {
				h.LOG.I.Println("file: ", entry.Name())
				err = h.parseEPUB(path)
				h.moveFile(path, err)
				if err != nil {
					h.LOG.W.Println(err)
				}
			}()
		case ext == ".zip":
			start := time.Now()
			new := !h.Hashes.ArchiveExists(entry.Name())
			h.LOG.I.Println("zip: ", entry.Name())
			err = h.parseZipEntry(path)
			h.moveFile(path, err)
			if err != nil {
				h.LOG.W.Println(err)
			}
			if new {
				h.LOG.S.Printf("%v elapsed for parsing %s ", time.Since(start), entry.Name())
			}
		default:
			h.LOG.D.Printf("file %s has not supported format \"%s\"\n", path, filepath.Ext(path))
			h.addFileToBookQueue(entry.Name(), "", hash.UnsupportedFormat)
			h.moveFile(path, err)
		}
	}
	return nil
}
