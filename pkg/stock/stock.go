package stock

import (
	"archive/zip"
	"fmt"
	"hash/crc32"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/vinser/flibgolite/pkg/config"
	"github.com/vinser/flibgolite/pkg/database"
	"github.com/vinser/flibgolite/pkg/fb2"
	"github.com/vinser/flibgolite/pkg/genres"
	"github.com/vinser/flibgolite/pkg/model"
	"github.com/vinser/flibgolite/pkg/parser"
	"github.com/vinser/flibgolite/pkg/rlog"
)

type Handler struct {
	CFG *config.Config
	DB  *database.DB
	GT  *genres.GenresTree
	LOG *rlog.Log
	SY  Sync
}

type Sync struct {
	WG         *sync.WaitGroup
	MaxThreads chan struct{}
	Stop       chan struct{}
}

// InitStockFolders()
func (h *Handler) InitStockFolders() {
	for _, stockDir := range []string{h.CFG.Library.STOCK_DIR, h.CFG.Library.NEW_DIR, h.CFG.Library.TRASH_DIR} {
		if err := os.MkdirAll(stockDir, 0775); err != nil {
			h.LOG.E.Printf("failed to create directory %s: %s", stockDir, err)
			log.Fatalf("failed to create directory %s: %s", stockDir, err)
		}
	}
}

// Reindex() - recr
func (h *Handler) Reindex() {
	db := h.DB
	db.DropDB()
	db.InitDB()
	start := time.Now()
	h.LOG.I.Println(">>> Book stock reindex started  >>>>>>>>>>>>>>>>>>>>>>>>>>>")
	h.ScanDir(h.CFG.Library.STOCK_DIR)
	finish := time.Now()
	h.LOG.I.Println("<<< Book stock reindex finished <<<<<<<<<<<<<<<<<<<<<<<<<<<")
	elapsed := finish.Sub(start)
	h.LOG.I.Println("Time elapsed: ", elapsed)
}

// Scan
func (h *Handler) ScanDir(stockDir string) error {
	// if !filepath.IsAbs(stockDir) {
	// 	workDir, _ := os.Getwd()
	// 	stockDir = filepath.Join(workDir, stockDir)
	// }
	d, err := os.Open(stockDir)
	if err != nil {
		return err
	}
	defer d.Close()
	entries, err := d.Readdir(-1)
	if err != nil {
		return err
	}
	h.SY.WG = &sync.WaitGroup{}
	h.SY.MaxThreads = make(chan struct{}, h.CFG.Database.MAX_SCAN_THREADS)
	for _, entry := range entries {
		path := filepath.Join(stockDir, entry.Name())
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		switch {
		case entry.Size() == 0:
			h.LOG.W.Printf("file %s from dir has size of zero\n", entry.Name())
			os.Rename(path, filepath.Join(h.CFG.Library.TRASH_DIR, entry.Name()))
		case entry.IsDir():
			h.LOG.W.Printf("subdirectory %s has been skipped\n ", path)
			// scanDir(false) // uncomment for recurse
		case ext == ".zip":
			h.LOG.I.Println("zip: ", entry.Name())
			h.SY.WG.Add(1)
			h.SY.MaxThreads <- struct{}{}
			go func() {
				err = h.processZip(path)
				h.moveFile(path, err)
				if err != nil {
					h.LOG.W.Println(err)
				}
			}()
		default:
			h.LOG.I.Println("file: ", entry.Name())
			err = h.processFile(path)
			h.moveFile(path, err)
			if err != nil {
				h.LOG.W.Println(err)
			}
		}
	}
	h.SY.WG.Wait()
	return nil
}

// Process single FB2 file
func (h *Handler) processFile(path string) error {
	crc32 := fileCRC32(path)
	fInfo, _ := os.Stat(path)
	if h.DB.IsInStock(fInfo.Name(), crc32) {
		return fmt.Errorf("file %s is in stock already and has been skipped", path)
	}
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %s", path, err)
	}
	defer f.Close()

	var p parser.Parser
	switch filepath.Ext(path) {
	case ".fb2":
		p, err = fb2.NewFB2(f)
		if err != nil {
			return fmt.Errorf("file %s has errors: %s", path, err)
		}
	default:
		return fmt.Errorf("file %s has unsupported format \"%s\"", path, filepath.Ext(path))
	}
	h.LOG.D.Println(p)
	book := &model.Book{
		File:     fInfo.Name(),
		CRC32:    crc32,
		Archive:  "",
		Size:     fInfo.Size(),
		Format:   p.GetFormat(),
		Title:    p.GetTitle(),
		Sort:     p.GetSort(),
		Year:     p.GetYear(),
		Plot:     p.GetPlot(),
		Cover:    p.GetCover(),
		Language: p.GetLanguage(),
		Authors:  p.GetAuthors(),
		Genres:   p.GetGenres(),
		Keywords: p.GetKeywords(),
		Serie:    p.GetSerie(),
		SerieNum: p.GetSerieNumber(),
		Updated:  time.Now().Unix(),
	}
	if !h.acceptLanguage(book.Language.Code) {
		return fmt.Errorf("publication language \"%s\" is configured as not accepted, file %s has been skipped", book.Language.Code, path)
	}
	h.GT.Refine(book)
	h.DB.NewBook(book)
	h.LOG.I.Printf("file %s has been added\n", path)
	return nil
}

// Process zip archive with FB2 files
func (h *Handler) processZip(zipPath string) error {
	defer h.SY.WG.Done()
	zr, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("incorrect zip archive %s", zipPath)
	}
	defer zr.Close()

	for _, file := range zr.File {
		h.LOG.D.Print(ZipEntryInfo(file))
		if filepath.Ext(file.Name) != ".fb2" {
			h.LOG.E.Printf("file %s from %s has non FB2 format\n", file.Name, filepath.Base(zipPath))
			continue
		}
		if h.DB.IsInStock(file.Name, file.CRC32) {
			h.LOG.W.Printf("file %s from %s is in stock already and has been skipped\n", file.Name, filepath.Base(zipPath))
			continue
		}
		if file.UncompressedSize == 0 {
			h.LOG.W.Printf("file %s from %s has size of zero\n", file.Name, filepath.Base(zipPath))
			continue
		}
		f, _ := file.Open()
		var p parser.Parser
		switch filepath.Ext(file.Name) {
		case ".fb2":
			p, err = fb2.NewFB2(f)
			if err != nil {
				h.LOG.E.Printf("file %s from %s has error: %s\n", file.Name, filepath.Base(zipPath), err.Error())
				f.Close()
				continue
			}
		default:
			h.LOG.W.Printf("file %s has unsupported format \"%s\"\n", file.Name, filepath.Ext(file.Name))
		}
		h.LOG.D.Println(p)
		book := &model.Book{
			File:     file.Name,
			CRC32:    file.CRC32,
			Archive:  filepath.Base(zipPath),
			Size:     int64(file.UncompressedSize),
			Format:   p.GetFormat(),
			Title:    p.GetTitle(),
			Sort:     p.GetSort(),
			Year:     p.GetYear(),
			Plot:     p.GetPlot(),
			Cover:    p.GetCover(),
			Language: p.GetLanguage(),
			Authors:  p.GetAuthors(),
			Genres:   p.GetGenres(),
			Keywords: p.GetKeywords(),
			Serie:    p.GetSerie(),
			SerieNum: p.GetSerieNumber(),
			Updated:  time.Now().Unix(),
		}
		if !h.acceptLanguage(book.Language.Code) {
			h.LOG.W.Printf("publication language \"%s\" is not accepted, file %s from %s has been skipped\n", book.Language.Code, file.Name, filepath.Base(zipPath))
			continue
		}
		h.GT.Refine(book)
		h.DB.NewBook(book)
		f.Close()
		h.LOG.I.Printf("file %s from %s has been added\n", file.Name, filepath.Base(zipPath))

		// runtime.Gosched()
	}
	<-h.SY.MaxThreads
	return nil
}

func (h *Handler) acceptLanguage(lang string) bool {
	return strings.Contains(h.CFG.Locales.ACCEPTED, lang)
}

func (h *Handler) moveFile(filePath string, err error) {
	if err != nil {
		os.Rename(filePath, filepath.Join(h.CFG.Library.TRASH_DIR, filepath.Base(filePath)))
		return
	}
	if filepath.Dir(filePath) == h.CFG.Library.STOCK_DIR {
		return
	}
	os.Rename(filePath, filepath.Join(h.CFG.Library.STOCK_DIR, filepath.Base(filePath)))
}

// fileCRC32 calculates file CRC32
func fileCRC32(filePath string) uint32 {
	fbytes, err := os.ReadFile(filePath)
	if err != nil {
		return 0
	}
	return crc32.ChecksumIEEE(fbytes)
}

//===============================
func ZipEntryInfo(e *zip.File) string {
	return "\n===========================================\n" +
		fmt.Sprintln("File               : ", e.Name) +
		fmt.Sprintln("NonUTF8            : ", e.NonUTF8) +
		fmt.Sprintln("Modified           : ", e.Modified) +
		fmt.Sprintln("CRC32              : ", e.CRC32) +
		fmt.Sprintln("UncompressedSize64 : ", e.UncompressedSize) +
		"===========================================\n"
}
