package opds

import (
	"io"
	"mime"
	"net/http"
	"strings"

	"github.com/vinser/flibgolite/internal/core/config"
	"github.com/vinser/flibgolite/internal/genres"
	"github.com/vinser/flibgolite/internal/rlog"
	"github.com/vinser/flibgolite/internal/store"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Handler struct {
	CFG *config.Config
	LOG *rlog.Log
	DB  *store.DB
	GT  *genres.GenresTree
	MP  map[string]*message.Printer
}

func init() {
	_ = mime.AddExtensionType(".mobi", "application/x-mobipocket-ebook")
	_ = mime.AddExtensionType(".epub", "application/epub+zip")
	_ = mime.AddExtensionType(".cbz", "application/x-cbz")
	_ = mime.AddExtensionType(".cbr", "application/x-cbr")
	_ = mime.AddExtensionType(".fb2", "application/fb2")
	_ = mime.AddExtensionType(".fb2.zip", "application/fb2+zip")   // Zipped fb2
	_ = mime.AddExtensionType(".fb2.epub", "application/epub+zip") // Converted from fb2
	_ = mime.AddExtensionType(".pdf", "application/pdf")           // Overwrite default mime type
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.LOG.I.Println(commentURL("Router", r))
	// switch r.URL.Path {
	switch strings.ReplaceAll(r.URL.Path, "//", "/") { // compensate PocketBook Reader search query error
	case "/favicon.ico":
		h.unloadFavicon(w)
	case "/opds":
		h.root(w, r)
	case "/opds/latest":
		h.latest(w, r)
	case "/opds/languages":
		h.languages(w, r)
	case "/opds/opensearch":
		h.openSerach(w, r)
	case "/opds/search":
		h.serach(w, r)
	case "/opds/authors":
		h.authors(w, r)
	case "/opds/genres":
		h.genres(w, r)
	case "/opds/series":
		h.series(w, r)
	case "/opds/books":
		h.books(w, r)
	case "/opds/covers":
		h.covers(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Bad request"}`)
		return
	}
}

func (h *Handler) getLanguage(r *http.Request) string {
	lang := r.FormValue("language")
	if lang == "" {
		lang = h.CFG.Locales.DEFAULT
	}
	if _, ok := h.CFG.Locales.Languages[lang]; ok {
		return lang
	}
	t, _, err := language.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
	if err != nil || len(t) == 0 {
		return h.CFG.Locales.DEFAULT
	}
	tag, _, _ := h.CFG.Matcher.Match(t...)
	base, _ := tag.Base()
	return base.String()
}
