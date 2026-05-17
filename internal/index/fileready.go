package index

import (
	"archive/zip"
	"fmt"
	"io/fs"
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
