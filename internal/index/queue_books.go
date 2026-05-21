package index

import (
	"os"
	"strings"
	"time"

	"github.com/vinser/flibgolite/internal/core/model"
	"github.com/vinser/flibgolite/internal/hash"
	"github.com/vinser/flibgolite/internal/store"
)

// ------------------------------------------------------------
// Helpers
// ------------------------------------------------------------

func isSQLiteBusy(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "SQLITE_BUSY") ||
		strings.Contains(msg, "database is locked") ||
		strings.Contains(msg, "busy")
}

// commitBatch commits the current batch and updates hash state.
// Retries whole batch on SQLITE_BUSY.
func (h *Handler) commitBatch(tx *store.TX, pending []model.Book, level string) error {
	err := tx.TxEnd()
	if err == nil {
		for i := range pending {
			h.Hashes.AddBook(&pending[i], level)
		}
		return nil
	}

	if isSQLiteBusy(err) {
		h.LOG.W.Println("Commit SQLITE_BUSY, retrying batch")

		newTx := h.DB.TxBegin()

		for _, b := range pending {
			st := h.Hashes.GetState(&b, level)
			switch st {
			case hash.Unique:
				newTx.NewBook(&b)
			default:
				newTx.RecordBookState(&b, st)
			}
		}

		if cerr := newTx.TxEnd(); cerr != nil {
			return cerr
		}

		for i := range pending {
			h.Hashes.AddBook(&pending[i], level)
		}

		return nil
	}

	return err
}

// retrySingleBook handles retry logic for a single failed book write.
func (h *Handler) retrySingleBook(tx *store.TX, book *model.Book, state hash.BookState, level string) error {
	if tx != nil {
		tx.TxEnd()
	}

	tx = h.DB.TxBegin()

	var err error
	switch state {
	case hash.Unique:
		err = tx.NewBook(book)
	default:
		err = tx.RecordBookState(book, state)
	}

	if err != nil {
		h.LOG.W.Printf("Retry failed for %s, marking as error", book.File)
		state = hash.FileHasErrors
		tx.RecordBookState(book, state)
	}

	cerr := tx.TxEnd()
	if cerr != nil {
		return cerr
	}

	h.Hashes.AddBook(book, level)

	return nil
}

// ------------------------------------------------------------
// Main worker
// ------------------------------------------------------------

func (h *Handler) ProcessBookQueue() {
	var (
		tx       *store.TX
		bookInTX int
		pending  []model.Book
		level    = h.CFG.Database.DEDUPLICATE_LEVEL
		maxInTx  = h.CFG.Database.MAX_BOOKS_IN_TX
	)

	defer func() {
		if tx != nil {
			tx.TxEnd()
		}
		h.StopDB <- struct{}{}
	}()

	for {
		select {

		// ---------- New book from queue -----------------------------
		case book := <-h.BookQueue:

			// Start new batch
			if bookInTX == 0 {
				pending = nil
				tx = h.DB.TxBegin()
			}

			// Determine state (read-only)
			state := h.Hashes.GetState(&book, level)

			// Try to write book into batch
			var err error
			switch state {
			case hash.Unique:
				err = tx.NewBook(&book)
			default:
				err = tx.RecordBookState(&book, state)
			}

			// ---------- Error wrighting a single book -------------------
			if err != nil {
				h.LOG.W.Printf("Error writing book %s: %v", book.File, err)

				if fatal := h.retrySingleBook(tx, &book, state, level); fatal != nil {
					h.LOG.E.Println("Fatal commit error:", fatal)
					os.Exit(1)
				}

				// Reset batch state
				bookInTX = 0
				continue
			}

			// ---------- Successful book wright --------------------------
			bookInTX++
			pending = append(pending, book)

			if book.Archive == "" {
				h.LOG.I.Printf("single file %s has been added\n", book.File)
			} else {
				h.LOG.I.Printf("file %s from %s has been added\n", book.File, book.Archive)
			}

			// ---------- Commit by limit ---------------------------------
			if bookInTX >= maxInTx {
				h.LOG.D.Printf("Commit by limit (%d books)", bookInTX)

				if fatal := h.commitBatch(tx, pending, level); fatal != nil {
					h.LOG.E.Println("Fatal commit error:", fatal)
					os.Exit(1)
				}

				bookInTX = 0
			}

		// ---------- Commit by timeout -------------------------------
		case <-time.After(time.Second):
			h.LOG.D.Printf("Book queue timeout")

			if bookInTX > 0 && tx != nil {
				h.LOG.D.Printf("Commit by timeout (%d books)", bookInTX)

				if fatal := h.commitBatch(tx, pending, level); fatal != nil {
					h.LOG.E.Println("Fatal commit error:", fatal)
					os.Exit(1)
				}

				bookInTX = 0
			}

		// ---------- Stop Signal -------------------------------------
		case <-h.StopDB:
			return
		}
	}
}
