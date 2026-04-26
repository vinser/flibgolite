package index

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/vinser/flibgolite/internal/core/config"
	"github.com/vinser/flibgolite/internal/core/model"
	"github.com/vinser/flibgolite/internal/genres"
	"github.com/vinser/flibgolite/internal/hash"
	"github.com/vinser/flibgolite/internal/parsers"
	"github.com/vinser/flibgolite/internal/parsers/fb2"
	"github.com/vinser/flibgolite/internal/rlog"
	"github.com/vinser/flibgolite/internal/store"
)

type Handler struct {
	CFG       *config.Config
	Hashes    *hash.BookHashes
	DB        *store.DB
	GT        *genres.GenresTree
	LOG       *rlog.Log
	ScanWG    sync.WaitGroup
	FileQueue chan File
	BookQueue chan model.Book
	StopScan  chan struct{}
	StopDB    chan struct{}
}

type File struct {
	Reader  io.ReadCloser
	Name    string
	CRC32   uint32
	Archive string
	Size    int64
}

// InitStockFolders()
func (h *Handler) InitStockFolders() {
	if err := os.MkdirAll(h.CFG.Library.STOCK_DIR, 0776); err != nil {
		log.Fatalf("failed to create Library STOCK_DIR directory %s: %s", h.CFG.Library.STOCK_DIR, err)
	}
	if len(h.CFG.Library.TRASH_DIR) > 0 {
		if err := os.MkdirAll(h.CFG.Library.TRASH_DIR, 0776); err != nil {
			log.Fatalf("failed to create Library TRASH_DIR directory %s: %s", h.CFG.Library.TRASH_DIR, err)
		}
	}
	if len(h.CFG.Library.NEW_DIR) > 0 {
		if err := os.MkdirAll(h.CFG.Library.NEW_DIR, 0776); err != nil {
			log.Fatalf("failed to create Library NEW_DIR directory %s: %s", h.CFG.Library.NEW_DIR, err)
		}
	}
}

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

// Scan
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

func (h *Handler) addFileToBookQueue(file, archive string, state hash.BookState) {
	h.BookQueue <- model.Book{
		File:    file,
		Archive: archive,
		Updated: int64(state),
	}
}

func (h *Handler) ParseFB2Queue() {
	for {
		select {
		case file := <-h.FileQueue:
			func() {
				f := file.Reader
				defer func() {
					f.Close()
					h.ScanWG.Done()
				}()
				var p parsers.Parser
				p, err := fb2.ParseFB2(f)
				if err != nil {
					h.addFileToBookQueue(file.Name, file.Archive, hash.FileHasErrors)
					h.LOG.D.Printf("file %s from %s has error: <%s> and has been skipped\n", file.Name, file.Archive, err.Error())
					return
				}
				language := p.GetLanguage()
				if !h.acceptLanguage(language.Code) {
					h.addFileToBookQueue(file.Name, file.Archive, hash.LanguageNotAccepted)
					h.LOG.D.Printf("publication language \"%s\" is not accepted, file %s from %s has been skipped\n", language.Code, file.Name, file.Archive)
					return
				}
				h.LOG.D.Println(p)
				book := &model.Book{
					File:     file.Name,
					CRC32:    file.CRC32,
					Archive:  file.Archive,
					Size:     file.Size,
					Format:   p.GetFormat(),
					Title:    p.GetTitle(),
					Sort:     p.GetSort(),
					Year:     p.GetYear(),
					Plot:     p.GetPlot(),
					Cover:    p.GetCover(),
					Language: language,
					Authors:  p.GetAuthors(),
					Genres:   p.GetGenres(),
					Keywords: p.GetKeywords(),
					Serie:    p.GetSerie(),
					SerieNum: p.GetSerieNumber(),
					Updated:  time.Now().UnixNano(),
				}
				h.GT.Refine(book)
				h.BookQueue <- *book
			}()

		case <-time.After(time.Second):
			h.LOG.D.Printf("File queue timeout")
		case <-h.StopScan:
			return
		}
	}

}

func (h *Handler) AddBooksToIndex() {
	tx := &store.TX{}
	defer func() {
		tx.TxEnd()
		h.StopDB <- struct{}{}
	}()
	bookInTX := 0
	for {
		select {
		case book := <-h.BookQueue:
			h.Hashes.Add(book.File, book.Archive)
			if bookInTX == 0 {
				tx = h.DB.TxBegin()
			}
			switch state := h.Hashes.GetState(&book, h.CFG.Database.DEDUPLICATE_LEVEL); state {
			case hash.Unique:
				tx.NewBook(&book)
			default:
				tx.RecordBookState(&book, state)
			}
			bookInTX++
			if book.Archive == "" {
				h.LOG.I.Printf("single file %s has been added\n", book.File)
			} else {
				h.LOG.I.Printf("file %s from %s has been added\n", book.File, book.Archive)
			}
			if bookInTX >= h.CFG.Database.MAX_BOOKS_IN_TX {
				tx.TxEnd()
				bookInTX = 0
			}
		case <-time.After(time.Second):
			h.LOG.D.Printf("Book queue timeout")
			if tx.Tx != nil {
				tx.TxEnd()
			}
			bookInTX = 0
		case <-h.StopDB:
			return
		}
	}
}

func (h *Handler) acceptLanguage(lang string) bool {
	if strings.Contains(h.CFG.ACCEPTED, "any") {
		return true
	}
	return strings.Contains(h.CFG.ACCEPTED, lang)
}

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

// ===============================
func ZipEntryInfo(e *zip.File) string {
	return "\n===========================================\n" +
		fmt.Sprintln("File               : ", e.Name) +
		fmt.Sprintln("NonUTF8            : ", e.NonUTF8) +
		fmt.Sprintln("Modified           : ", e.Modified) +
		fmt.Sprintln("CRC32              : ", e.CRC32) +
		fmt.Sprintln("UncompressedSize64 : ", e.UncompressedSize64) +
		"===========================================\n"
}
