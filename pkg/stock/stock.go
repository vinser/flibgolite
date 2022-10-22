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
	"github.com/vinser/flibgolite/pkg/epub"
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
	if err := os.MkdirAll(h.CFG.Library.STOCK_DIR, 0666); err != nil {
		log.Fatalf("failed to create Library STOCK_DIR directory %s: %s", h.CFG.Library.STOCK_DIR, err)
	}
	if err := os.MkdirAll(h.CFG.Library.TRASH_DIR, 0666); err != nil {
		log.Fatalf("failed to create Library TRASH_DIR directory %s: %s", h.CFG.Library.TRASH_DIR, err)
	}
	if len(h.CFG.Library.NEW_DIR) > 0 {
		if err := os.MkdirAll(h.CFG.Library.NEW_DIR, 0666); err != nil {
			log.Fatalf("failed to create Library NEW_DIR directory %s: %s", h.CFG.Library.NEW_DIR, err)
		}
	}
}

// Reindex() - recreate book stock database
func (h *Handler) Reindex() {
	db := h.DB
	db.DropDB()
	db.InitDB()
	start := time.Now()
	h.LOG.S.Println(">>> Book stock reindex started  >>>>>>>>>>>>>>>>>>>>>>>>>>>")
	h.ScanDir(h.CFG.Library.STOCK_DIR)
	finish := time.Now()
	h.LOG.S.Println("<<< Book stock reindex finished <<<<<<<<<<<<<<<<<<<<<<<<<<<")
	elapsed := finish.Sub(start)
	h.LOG.S.Println("Time elapsed: ", elapsed)
}

// Scan
func (h *Handler) ScanDir(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	entries, err := d.Readdir(-1)
	if err != nil {
		return err
	}
	h.LOG.I.Printf("scanning folder %s for new books.../n", dir)
	h.SY.WG = &sync.WaitGroup{}
	h.SY.MaxThreads = make(chan struct{}, h.CFG.Database.MAX_SCAN_THREADS)
	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		switch {
		case entry.Size() == 0:
			h.LOG.W.Printf("file %s from dir has size of zero\n", entry.Name())
			os.Rename(path, filepath.Join(h.CFG.Library.TRASH_DIR, entry.Name()))
		case entry.IsDir():
			h.LOG.W.Printf("subdirectory %s has been skipped\n ", path)
		case ext == ".fb2":
			h.LOG.D.Println("file: ", entry.Name())
			err = h.indexFB2File(path)
			h.moveFile(path, err)
			if err != nil {
				h.LOG.W.Println(err)
			}
		case ext == ".epub":
			h.LOG.D.Println("file: ", entry.Name())
			err = h.indexEPUBFile(path)
			h.moveFile(path, err)
			if err != nil {
				h.LOG.W.Println(err)
			}
		case ext == ".zip":
			h.LOG.D.Println("zip: ", entry.Name())
			h.SY.WG.Add(1)
			h.SY.MaxThreads <- struct{}{}
			go func() {
				err = h.indexFB2Zip(path)
				h.moveFile(path, err)
				if err != nil {
					h.LOG.W.Println(err)
				}
			}()
		default:
			h.LOG.E.Printf("file %s has not supported format \"%s\"\n", path, filepath.Ext(path))
			h.moveFile(path, err)
		}
	}
	h.SY.WG.Wait()
	return nil
}

// Process single FB2 file and add it to book stock index
func (h *Handler) indexFB2File(FB2Path string) error {
	crc32 := fileCRC32(FB2Path)
	fInfo, _ := os.Stat(FB2Path)
	if h.DB.IsFileInStock(fInfo.Name(), crc32) {
		if len(h.CFG.Library.NEW_DIR) > 0 {
			return fmt.Errorf("file %s is in stock already and has been skipped", FB2Path)
		}
		h.LOG.D.Printf("file %s is in stock already and has been skipped\n", FB2Path)
		return nil
	}
	f, err := os.Open(FB2Path)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %s", FB2Path, err)
	}
	defer f.Close()

	var p parser.Parser
	p, err = fb2.NewFB2(f)
	if err != nil {
		return fmt.Errorf("file %s has errors: %s", FB2Path, err)
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
		return fmt.Errorf("publication language \"%s\" is configured as not accepted, file %s has been skipped", book.Language.Code, FB2Path)
	}
	h.GT.Refine(book)
	h.DB.NewBook(book)
	h.LOG.I.Printf("file %s has been added\n", FB2Path)
	return nil
}

// Process single EPUB file and add it to book stock index
func (h *Handler) indexEPUBFile(EPUBPath string) error {
	crc32 := fileCRC32(EPUBPath)
	fInfo, _ := os.Stat(EPUBPath)
	if h.DB.IsFileInStock(fInfo.Name(), crc32) {
		if len(h.CFG.Library.NEW_DIR) > 0 {
			return fmt.Errorf("file %s is in stock already and has been skipped", EPUBPath)
		}
		h.LOG.D.Printf("file %s is in stock already and has been skipped\n", EPUBPath)
		return nil
	}

	zr, err := zip.OpenReader(EPUBPath)
	if err != nil {
		return fmt.Errorf("incorrect zip archive %s", EPUBPath)
	}
	defer zr.Close()

	var p parser.Parser
	zPath, err := epub.GetOPFPath(zr)
	if err != nil {
		return fmt.Errorf("file %s has errors: %s", EPUBPath, err)
	}
	p, err = epub.NewOPF(zr, zPath)
	if err != nil {
		return fmt.Errorf("file %s has errors: %s", EPUBPath, err)
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
		return fmt.Errorf("publication language \"%s\" is configured as not accepted, file %s has been skipped", book.Language.Code, EPUBPath)
	}
	h.GT.Refine(book)
	h.DB.NewBook(book)
	h.LOG.I.Printf("file %s has been added\n", EPUBPath)
	return nil
}

// Process zip archive with FB2 files and add them to book stock index
func (h *Handler) indexFB2Zip(zipPath string) error {
	defer h.SY.WG.Done()
	defer func() {
		<-h.SY.MaxThreads
	}()
	if h.DB.IsArchiveInStock(filepath.Base(zipPath)) {
		if len(h.CFG.Library.NEW_DIR) > 0 {
			return fmt.Errorf("archive %s is in stock already and has been skipped", zipPath)
		}
		h.LOG.D.Printf("archive %s is in stock already and has been skipped\n", zipPath)
		return nil
	}
	zr, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("incorrect zip archive %s", zipPath)
	}
	defer zr.Close()

	for _, file := range zr.File {
		h.LOG.D.Print(ZipEntryInfo(file))
		if h.DB.IsFileInStock(file.Name, file.CRC32) {
			h.LOG.D.Printf("file %s from %s is in stock already and has been skipped\n", file.Name, filepath.Base(zipPath))
			continue
		}
		if file.UncompressedSize == 0 {
			h.LOG.W.Printf("file %s from %s has size of zero and has been skipped\n", file.Name, filepath.Base(zipPath))
			continue
		}
		f, _ := file.Open()
		var p parser.Parser
		switch filepath.Ext(file.Name) {
		case ".fb2":
			p, err = fb2.NewFB2(f)
			if err != nil {
				h.LOG.E.Printf("file %s from %s has error: <%s> and has been skipped\n", file.Name, filepath.Base(zipPath), err.Error())
				f.Close()
				continue
			}
		default:
			h.LOG.D.Printf("file %s from %s has unsupported format \"%s\" and has been skipped\n", file.Name, filepath.Base(zipPath), filepath.Ext(file.Name))
			continue
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
	// <-h.SY.MaxThreads
	h.LOG.I.Printf("archive %s has been indexed\n", zipPath)
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

// ===============================
func ZipEntryInfo(e *zip.File) string {
	return "\n===========================================\n" +
		fmt.Sprintln("File               : ", e.Name) +
		fmt.Sprintln("NonUTF8            : ", e.NonUTF8) +
		fmt.Sprintln("Modified           : ", e.Modified) +
		fmt.Sprintln("CRC32              : ", e.CRC32) +
		fmt.Sprintln("UncompressedSize64 : ", e.UncompressedSize) +
		"===========================================\n"
}
