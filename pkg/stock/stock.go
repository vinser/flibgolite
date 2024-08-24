package stock

import (
	"archive/zip"
	"fmt"
	"hash/crc32"
	"io"
	"io/fs"
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
	CFG   *config.Config
	DB    *database.DB
	TX    *database.TX
	GT    *genres.GenresTree
	LOG   *rlog.Log
	WG    sync.WaitGroup
	Queue chan File
	Stop  chan struct{}
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
	if err := os.MkdirAll(h.CFG.Library.TRASH_DIR, 0776); err != nil {
		log.Fatalf("failed to create Library TRASH_DIR directory %s: %s", h.CFG.Library.TRASH_DIR, err)
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
		err = h.DB.NotInStock(info.Name())
		if err != nil {
			return "", "", err
		}
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
					os.Rename(path, filepath.Join(h.CFG.Library.TRASH_DIR, info.Name()))
					return "", "", fmt.Errorf("file %s has size of zero", path)
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
	return "", "", fmt.Errorf("file %s is not a regular file", path)
}

// Scan
func (h *Handler) ScanDir(dir string, queue chan<- model.Book) error {
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
				h.LOG.D.Println("file: ", entry.Name())
				err = h.indexFB2File(path, queue)
				h.moveFile(path, err)
				if err != nil {
					h.LOG.W.Println(err)
				}
			}()
		case ext == ".epub":
			go func() {
				h.LOG.D.Println("file: ", entry.Name())
				err = h.indexEPUBFile(path, queue)
				h.moveFile(path, err)
				if err != nil {
					h.LOG.W.Println(err)
				}
			}()
		case ext == ".zip":
			start := time.Now()
			h.LOG.D.Println("zip: ", entry.Name())
			err = h.indexFB2Zip(path)
			h.moveFile(path, err)
			if err != nil {
				h.LOG.W.Println(err)
			}
			h.LOG.S.Printf("%v elapsed for parsing %s ", time.Since(start), entry.Name())
		default:
			h.LOG.E.Printf("file %s has not supported format \"%s\"\n", path, filepath.Ext(path))
			h.moveFile(path, err)
		}
	}
	return nil
}

// Process single FB2 file and add it to book stock index
func (h *Handler) indexFB2File(FB2Path string, queue chan<- model.Book) error {
	fInfo, _ := os.Stat(FB2Path)
	f, err := os.Open(FB2Path)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %s", FB2Path, err)
	}
	defer f.Close()

	var p parser.Parser
	p, err = fb2.ParseFB2(f)
	if err != nil {
		return fmt.Errorf("file %s has errors: %s", FB2Path, err)
	}
	h.LOG.D.Println(p)
	language := p.GetLanguage()
	if !h.acceptLanguage(language.Code) {
		return fmt.Errorf("publication language \"%s\" is configured as not accepted, file %s has been skipped", language.Code, FB2Path)
	}
	crc32 := fileCRC32(FB2Path)
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
		Language: language,
		Authors:  p.GetAuthors(),
		Genres:   p.GetGenres(),
		Keywords: p.GetKeywords(),
		Serie:    p.GetSerie(),
		SerieNum: p.GetSerieNumber(),
		Updated:  time.Now().Unix(),
	}
	h.GT.Refine(book)
	queue <- *book
	return nil
}

// Process single EPUB file and add it to book stock index
func (h *Handler) indexEPUBFile(EPUBPath string, queue chan<- model.Book) error {
	crc32 := fileCRC32(EPUBPath)
	fInfo, _ := os.Stat(EPUBPath)
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
	language := p.GetLanguage()
	if !h.acceptLanguage(language.Code) {
		return fmt.Errorf("publication language \"%s\" is configured as not accepted, file %s has been skipped", language.Code, EPUBPath)
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
		Language: language,
		Authors:  p.GetAuthors(),
		Genres:   p.GetGenres(),
		Keywords: p.GetKeywords(),
		Serie:    p.GetSerie(),
		SerieNum: p.GetSerieNumber(),
		Updated:  time.Now().Unix(),
	}
	h.GT.Refine(book)
	queue <- *book
	return nil
}

// Process zip archive with FB2 files and add them to book stock index
func (h *Handler) indexFB2Zip(zipPath string) error {
	h.LOG.D.Printf("archive %s indexing has been started\n", zipPath)
	zr, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("incorrect zip archive %s: %s", zipPath, err)
	}
	defer func() {
		zr.Close()
		h.LOG.D.Printf("archive %s indexing has been finished\n", zipPath)
	}()
	for _, file := range zr.File {
		h.LOG.D.Print(ZipEntryInfo(file))
		if file.UncompressedSize64 == 0 {
			h.LOG.I.Printf("file %s from %s has size of zero and has been skipped\n", file.Name, filepath.Base(zipPath))
			continue
		}
		if filepath.Ext(file.Name) != ".fb2" {
			h.LOG.I.Printf("file %s from %s has unsupported format \"%s\" and has been skipped\n", file.Name, filepath.Base(zipPath), filepath.Ext(file.Name))
			continue

		}

		h.WG.Add(1)
		f := &File{
			Reader: func() io.ReadCloser {
				r, _ := file.Open()
				return r
			}(),
			Name:    file.Name,
			CRC32:   file.CRC32,
			Archive: filepath.Base(zipPath),
			Size:    int64(file.UncompressedSize64),
		}
		h.Queue <- *f

	}
	h.WG.Wait()
	return nil
}

func (h *Handler) ParseFB2Queue(queue chan<- model.Book) {
	for {
		select {
		case file := <-h.Queue:
			func() {
				f := file.Reader
				defer func() {
					f.Close()
					h.WG.Done()
				}()
				var p parser.Parser
				p, err := fb2.ParseFB2(f)
				if err != nil {
					h.LOG.E.Printf("file %s from %s has error: <%s> and has been skipped\n", file.Name, file.Archive, err.Error())
					return
				}
				language := p.GetLanguage()
				if !h.acceptLanguage(language.Code) {
					h.LOG.I.Printf("publication language \"%s\" is not accepted, file %s from %s has been skipped\n", language.Code, file.Name, file.Archive)
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
					Updated:  time.Now().Unix(),
				}
				h.GT.Refine(book)
				queue <- *book
			}()

		case <-time.After(time.Second):
			h.LOG.D.Printf("File queue timeout")
		case <-h.Stop:
			return
		}
	}

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
		fmt.Sprintln("UncompressedSize64 : ", e.UncompressedSize64) +
		"===========================================\n"
}
