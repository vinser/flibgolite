package index

import (
	"time"

	"github.com/vinser/flibgolite/internal/hash"
	"github.com/vinser/flibgolite/internal/parsers"
	"github.com/vinser/flibgolite/internal/parsers/fb2"
	"github.com/vinser/flibgolite/internal/store"
)

// ParseFB2Queue processes files from the file queue
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

// AddBooksToIndex processes books from the book queue and adds them to the database
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
				err := tx.NewBook(&book)
				if err != nil {
					h.LOG.W.Println("Error adding book to database:", err)
					return
				}
			default:
				err := tx.RecordBookState(&book, state)
				if err != nil {
					h.LOG.W.Printf("Error recording book state for %s: %v", book.File, err)
					return
				}
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
