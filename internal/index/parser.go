package index

import (
	"archive/zip"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/vinser/flibgolite/internal/core/model"
	"github.com/vinser/flibgolite/internal/hash"
	"github.com/vinser/flibgolite/internal/parsers"
	"github.com/vinser/flibgolite/internal/parsers/epub"
	"github.com/vinser/flibgolite/internal/parsers/fb2"
)

// parseFB2 processes a single FB2 file and adds it to book stock index
func (h *Handler) parseFB2(FB2Path string) error {
	fInfo, _ := os.Stat(FB2Path)
	file := fInfo.Name()
	if h.Hashes.FileExists(file, "") {
		h.LOG.D.Printf("file %s is in stock already and has been skipped", file)
		return nil
	}
	f, err := os.Open(FB2Path)
	if err != nil {
		h.addFileToBookQueue(file, "", hash.FileOpenFailed)
		return fmt.Errorf("failed to open file %s: %s", FB2Path, err)
	}
	defer f.Close()

	var p parsers.Parser
	p, err = fb2.ParseFB2(f)
	if err != nil {
		h.addFileToBookQueue(file, "", hash.FileHasErrors)
		return fmt.Errorf("file %s has errors: %s", file, err)
	}
	h.LOG.D.Println(p)
	language := p.GetLanguage()
	if !h.acceptLanguage(language.Code) {
		h.addFileToBookQueue(file, "", hash.LanguageNotAccepted)
		return fmt.Errorf("publication language \"%s\" is configured as not accepted, file %s has been skipped", language.Code, file)
	}
	book := &model.Book{
		File:     file,
		CRC32:    fileCRC32(FB2Path),
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
		Updated:  time.Now().UnixNano(),
	}
	h.GT.Refine(book)
	h.BookQueue <- *book
	return nil
}

// parseEPUB processes a single EPUB file and adds it to book stock index
func (h *Handler) parseEPUB(EPUBPath string) error {
	fInfo, _ := os.Stat(EPUBPath)
	file := fInfo.Name()
	if h.Hashes.FileExists(file, "") {
		h.LOG.D.Printf("file %s is in stock already and has been skipped", file)
		return nil
	}
	zr, err := zip.OpenReader(EPUBPath)
	if err != nil {
		h.addFileToBookQueue(fInfo.Name(), "", hash.BadArchive)
		return fmt.Errorf("incorrect zip archive %s", file)
	}
	defer zr.Close()

	var p parsers.Parser
	zPath, err := epub.GetOPFPath(zr)
	if err != nil {
		h.addFileToBookQueue(fInfo.Name(), "", hash.FileHasErrors)
		return fmt.Errorf("file %s has errors: %s", file, err)
	}
	p, err = epub.NewOPF(zr, zPath)
	if err != nil {
		h.addFileToBookQueue(fInfo.Name(), "", hash.FileHasErrors)
		return fmt.Errorf("file %s has errors: %s", file, err)
	}
	language := p.GetLanguage()
	if !h.acceptLanguage(language.Code) {
		h.addFileToBookQueue(fInfo.Name(), "", hash.LanguageNotAccepted)
		return fmt.Errorf("publication language \"%s\" is configured as not accepted, file %s has been skipped", language.Code, file)
	}
	h.LOG.D.Println(p)
	book := &model.Book{
		File:     fInfo.Name(),
		CRC32:    fileCRC32(EPUBPath),
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
		Updated:  time.Now().UnixNano(),
	}
	h.GT.Refine(book)
	h.BookQueue <- *book
	return nil
}

// parseZipEntry processes a zip archive with FB2 files and adds them to book stock index
func (h *Handler) parseZipEntry(zipPath string) error {
	h.LOG.D.Printf("archive %s indexing has been started\n", zipPath)
	zr, err := zip.OpenReader(zipPath)
	if err != nil {
		h.addFileToBookQueue("", filepath.Base(zipPath), hash.BadArchive)
		return fmt.Errorf("incorrect zip archive %s: %s", zipPath, err)
	}
	defer func() {
		zr.Close()
		h.LOG.D.Printf("archive %s indexing has been finished\n", zipPath)
	}()
	for _, file := range zr.File {
		h.LOG.D.Print(ZipEntryInfo(file))

		if h.Hashes.FileExists(filepath.Base(file.Name), filepath.Base(zipPath)) {
			h.LOG.D.Printf("file %s from %s is in stock already and has been skipped", filepath.Base(file.Name), filepath.Base(zipPath))
			continue
		}

		if file.UncompressedSize64 == 0 {
			h.addFileToBookQueue(filepath.Base(file.Name), filepath.Base(zipPath), hash.FileIsEmpty)
			h.LOG.D.Printf("file %s from %s has size of zero and has been skipped\n", file.Name, filepath.Base(zipPath))
			continue
		}
		if filepath.Ext(file.Name) != ".fb2" {
			h.addFileToBookQueue(filepath.Base(file.Name), filepath.Base(zipPath), hash.UnsupportedFormat)
			h.LOG.D.Printf("file %s from %s has unsupported format \"%s\" and has been skipped\n", file.Name, filepath.Base(zipPath), filepath.Ext(file.Name))
			continue

		}

		h.ScanWG.Add(1)
		f := &File{
			Reader: func() io.ReadCloser {
				r, _ := file.Open()
				return r
			}(),
			Name:    filepath.Base(file.Name),
			CRC32:   file.CRC32,
			Archive: filepath.Base(zipPath),
			Size:    int64(file.UncompressedSize64),
		}
		h.FileQueue <- *f

	}
	h.ScanWG.Wait()
	return nil
}

// fileCRC32 calculates file CRC32
func fileCRC32(filePath string) uint32 {
	fbytes, err := os.ReadFile(filePath)
	if err != nil {
		return 0
	}
	return crc32.ChecksumIEEE(fbytes)
}
