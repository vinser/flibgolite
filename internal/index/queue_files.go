package index

import (
	"time"

	"github.com/vinser/flibgolite/internal/hash"
	"github.com/vinser/flibgolite/internal/parsers"
	"github.com/vinser/flibgolite/internal/parsers/fb2"
)

// ProcessFileQueue processes files from the file queue
func (h *Handler) ProcessFileQueue() {
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
				if err := h.processLanguage(p, file.Name, file.Archive); err != nil {
					return
				}
				h.LOG.D.Println(p)
				book := h.createBookFromParser(p, file.Name, file.Archive, file.Size, file.CRC32)
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
