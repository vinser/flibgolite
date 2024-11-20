package opds

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"math"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/nfnt/resize"
	"github.com/vinser/flibgolite/pkg/config"
	cfb2 "github.com/vinser/flibgolite/pkg/conv/fb2"
	"github.com/vinser/flibgolite/pkg/database"
	"github.com/vinser/flibgolite/pkg/epub"
	"github.com/vinser/flibgolite/pkg/fb2"
	"github.com/vinser/flibgolite/pkg/genres"
	"github.com/vinser/flibgolite/pkg/model"
	"github.com/vinser/flibgolite/pkg/rlog"
	"github.com/vinser/u8xml"

	"golang.org/x/text/cases"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
	"golang.org/x/text/message"

	_ "image/gif"
	_ "image/png"
)

type Handler struct {
	CFG *config.Config
	LOG *rlog.Log
	DB  *database.DB
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
	_ = mime.AddExtensionType(".fb2.zip", "application/fb2+zip")   // Zipped fb2
	_ = mime.AddExtensionType(".fb2.epub", "application/epub+zip") // Converted from fb2
	_ = mime.AddExtensionType(".pdf", "application/pdf")           // Overwrite default mime type
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.LOG.I.Println(commentURL("Router", r))
	// switch r.URL.Path {
	switch strings.ReplaceAll(r.URL.Path, "//", "/") { // compensate PocketBook Reader search query error
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

// Root
func (h *Handler) root(w http.ResponseWriter, r *http.Request) {
	lang := h.getLanguage(r)
	selfHref := fmt.Sprintf("/opds?language=%s", lang)
	f := NewFeed(h.CFG.OPDS.TITLE, "", selfHref)
	searchLink := &Link{Rel: FeedSearchLinkRel, Href: fmt.Sprintf("/opds/search?language=%s&q={searchTerms}", lang), Type: "application/atom+xml"}
	f.Link = append(f.Link, *searchLink)
	searchDescLink := &Link{Rel: FeedSearchLinkRel, Href: fmt.Sprintf("/opds/opensearch?language=%s", lang), Type: FeedSearchDescriptionLinkType, Title: "Search on catalog"}
	f.Link = append(f.Link, *searchDescLink)
	f.Entry = []*Entry{
		{
			Title:   h.MP[lang].Sprintf("Latest"),
			ID:      "authors",
			Updated: f.Time(time.Now()),
			Links: []Link{
				{
					Rel:  FeedSubsectionLinkRel,
					Href: fmt.Sprintf("/opds/latest?language=%s", lang),
					Type: FeedNavigationLinkType,
				},
			},
			Content: &Content{
				Type:    FeedTextContentType,
				Content: h.MP[lang].Sprintf("Browse the latest books received"),
			},
		},
		{
			Title:   h.MP[lang].Sprintf("Book Authors"),
			ID:      "authors",
			Updated: f.Time(time.Now()),
			Links: []Link{
				{
					Rel:  FeedSubsectionLinkRel,
					Href: fmt.Sprintf("/opds/authors?language=%s", lang),
					Type: FeedNavigationLinkType,
				},
			},
			Content: &Content{
				Type:    FeedTextContentType,
				Content: h.MP[lang].Sprintf("Browse books by author"),
			},
		},
		{
			Title:   h.MP[lang].Sprintf("Book Series"),
			ID:      "series",
			Updated: f.Time(time.Now()),
			Links: []Link{
				{
					Rel:  FeedSubsectionLinkRel,
					Href: fmt.Sprintf("/opds/series?language=%s", lang),
					Type: FeedNavigationLinkType,
				},
			},
			Content: &Content{
				Type:    FeedTextContentType,
				Content: h.MP[lang].Sprintf("Browse books by series"),
			},
		},
		{
			Title:   h.MP[lang].Sprintf("Book Genres"),
			ID:      "genres",
			Updated: f.Time(time.Now()),
			Links: []Link{
				{
					Rel:  FeedSubsectionLinkRel,
					Href: fmt.Sprintf("/opds/genres?language=%s", lang),
					Type: FeedNavigationLinkType,
				},
			},
			Content: &Content{
				Type:    FeedTextContentType,
				Content: h.MP[lang].Sprintf("Browse books by genre"),
			},
		},
		{
			Title:   h.MP[lang].Sprintf("Book Languages"),
			ID:      "languages",
			Updated: f.Time(time.Now()),
			Links: []Link{
				{
					Rel:  FeedSubsectionLinkRel,
					Href: fmt.Sprintf("/opds/languages?language=%s", lang),
					Type: FeedNavigationLinkType,
				},
			},
			Content: &Content{
				Type:    FeedTextContentType,
				Content: h.MP[lang].Sprintf("Language selection"),
			},
		},
	}
	//
	writeFeed(w, http.StatusOK, *f)
}

// Latest
func (h *Handler) latest(w http.ResponseWriter, r *http.Request) {
	lang := h.getLanguage(r)
	h.LOG.D.Println(commentURL("Latest", r))
	selfHref := ""
	bc := h.DB.LatestBooksCount(h.CFG.OPDS.LATEST_DAYS)

	switch {
	case bc != 0: // show books
		page, err := strconv.Atoi(r.FormValue("page"))
		if err != nil {
			page = 1
		}
		offset := (page - 1) * h.CFG.OPDS.PAGE_SIZE
		books := h.DB.PageLatestBooks(h.CFG.OPDS.LATEST_DAYS, h.CFG.OPDS.PAGE_SIZE+1, offset)
		selfHref = fmt.Sprintf("/opds/latest?language=%s&page=%d", lang, page)
		f := NewFeed(h.MP[lang].Sprintf("Found titles - %d", bc), "", selfHref)
		if len(books) > h.CFG.OPDS.PAGE_SIZE {
			nextRef := fmt.Sprintf("/opds/latest?language=%s&page=%d", lang, page+1)
			nextLink := &Link{Rel: FeedNextLinkRel, Href: nextRef, Type: FeedNavigationLinkType}
			f.Link = append(f.Link, *nextLink)
			books = books[:h.CFG.OPDS.PAGE_SIZE-1]
		}
		if int(bc) > h.CFG.OPDS.PAGE_SIZE {
			if page > 1 {
				firstRef := fmt.Sprintf("/opds/latest?language=%s&page=1", lang)
				firstLink := &Link{Rel: FeedFirstLinkRel, Href: firstRef, Type: FeedNavigationLinkType}
				f.Link = append(f.Link, *firstLink)

				prevRef := fmt.Sprintf("/opds/latest?language=%s&page=%d", lang, page-1)
				prevLink := &Link{Rel: FeedPrevLinkRel, Href: prevRef, Type: FeedNavigationLinkType}
				f.Link = append(f.Link, *prevLink)
			}
			lastPage := int(math.Ceil(float64(bc) / float64(h.CFG.OPDS.PAGE_SIZE)))
			if page < lastPage {
				lastRef := fmt.Sprintf("/opds/latest?language=%s&page=%d", lang, lastPage)
				lastLink := &Link{Rel: FeedLastLinkRel, Href: lastRef, Type: FeedNavigationLinkType}
				f.Link = append(f.Link, *lastLink)
			}
		}

		h.feedBookEntries(r, books, f)
		writeFeed(w, http.StatusOK, *f)
	default:
		selfHref = fmt.Sprintf("/opds/latest?language=%s", lang)
		f := NewFeed(h.MP[lang].Sprintf("Nothing found"), "", selfHref)
		writeFeed(w, http.StatusOK, *f)
	}
}

// Languages
func (h *Handler) languages(w http.ResponseWriter, r *http.Request) {
	lang := h.getLanguage(r)
	selfHref := fmt.Sprintf("/opds/languages?language=%s", lang)
	f := NewFeed(h.MP[lang].Sprintf("Language selection"), "", selfHref)
	ordered := []string{}
	for o := range h.CFG.Locales.Languages {
		ordered = append(ordered, o)
	}
	sort.Strings(ordered)
	for _, v := range ordered {
		langTitle := cases.Title(h.CFG.Locales.Languages[v].Tag)
		langName := langTitle.String(display.Self.Name(h.CFG.Locales.Languages[v].Tag))
		langBookTotal := h.DB.CountLanguageBooks(v)
		entry := &Entry{
			Title:   langName,
			ID:      "/opds/language=" + v,
			Updated: f.Time(time.Now()),
			Links: []Link{
				{
					Rel:  FeedSubsectionLinkRel,
					Href: fmt.Sprintf("/opds?language=%s", v),
					Type: FeedNavigationLinkType,
				},
			},
			Content: &Content{
				Type:    FeedTextContentType,
				Content: h.MP[lang].Sprintf("Total books - %d", langBookTotal),
			},
		}
		f.Entry = append(f.Entry, entry)
	}
	writeFeed(w, http.StatusOK, *f)
}

// OpenSearch description document
func (h *Handler) openSerach(w http.ResponseWriter, r *http.Request) {
	lang := h.getLanguage(r)
	data :=
		`
<OpenSearchDescription xmlns="http://a9.com/-/spec/opensearch/1.1/">
<ShortName>` + h.CFG.OPDS.TITLE + `</ShortName>
<Description>Search on catalog</Description>
<InputEncoding>UTF-8</InputEncoding>
<OutputEncoding>UTF-8</OutputEncoding>
<Url type="application/atom+xml" template=` + fmt.Sprintf("/opds/search?language=%s&q={searchTerms}", lang) + `/>
</OpenSearchDescription>	
`
	s := fmt.Sprintf("%s%s", xml.Header, data)
	w.Header().Add("Content-Type", "application/atom+xml")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, s)

}

// Search
func (h *Handler) serach(w http.ResponseWriter, r *http.Request) {
	lang := h.getLanguage(r)
	h.LOG.D.Println(commentURL("Search", r))
	selfHref := ""
	queryString := ""
	var ac, bc int64
	switch {
	case r.FormValue("q") != "":
		queryString = r.FormValue("q")
		if utf8.RuneCountInString(queryString) < 3 {
			return
		}
		bc = h.DB.SearchBooksCount(queryString)
		ac = h.DB.SearchAuthorsCount(queryString)
	case r.FormValue("book") != "":
		queryString = r.FormValue("book")
		bc = h.DB.SearchBooksCount(queryString)
	case r.FormValue("author") != "":
		queryString = r.FormValue("author")
		ac = h.DB.SearchAuthorsCount(queryString)
	}

	switch {
	case (ac != 0 && bc != 0):
		selfHref = fmt.Sprintf("/opds/search?language=%s&q={searchTerms}", lang)
		f := NewFeed(h.MP[lang].Sprintf("Choose from the found ones"), "", selfHref)
		f.Entry = []*Entry{
			{
				Title:   h.MP[lang].Sprintf("Titles"),
				ID:      fmt.Sprintf("/opds/search/book=%s", queryString),
				Updated: f.Time(time.Now()),
				Links: []Link{
					{
						Rel:  FeedSubsectionLinkRel,
						Href: fmt.Sprintf("/opds/search?language=%s&book=%s", lang, queryString),
						Type: FeedNavigationLinkType,
					},
				},
				Content: &Content{
					Type:    FeedTextContentType,
					Content: h.MP[lang].Sprintf("Found titles - %d", bc),
				},
			},
			{
				Title:   h.MP[lang].Sprintf("Authors"),
				ID:      fmt.Sprintf("/opds/search/author=%s", queryString),
				Updated: f.Time(time.Now()),
				Links: []Link{
					{
						Rel:  FeedSubsectionLinkRel,
						Href: fmt.Sprintf("/opds/search?language=%s&author=%s", lang, queryString),
						Type: FeedNavigationLinkType,
					},
				},
				Content: &Content{
					Type:    FeedTextContentType,
					Content: h.MP[lang].Sprintf("Found authors - %d", ac),
				},
			},
		}
		writeFeed(w, http.StatusOK, *f)
	case ac == 0 && bc != 0: // show books
		page, err := strconv.Atoi(r.FormValue("page"))
		if err != nil {
			page = 1
		}
		offset := (page - 1) * h.CFG.OPDS.PAGE_SIZE
		books := h.DB.PageFoundBooks(queryString, h.CFG.OPDS.PAGE_SIZE+1, offset)
		selfHref = fmt.Sprintf("/opds/search?language=%s&book=%s&page=%d", lang, queryString, page)
		f := NewFeed(h.MP[lang].Sprintf("Found titles - %d", bc), "", selfHref)
		if len(books) > h.CFG.OPDS.PAGE_SIZE {
			nextRef := fmt.Sprintf("/opds/search?language=%s&book=%s&page=%d", lang, queryString, page+1)
			nextLink := &Link{Rel: FeedNextLinkRel, Href: nextRef, Type: FeedNavigationLinkType}
			f.Link = append(f.Link, *nextLink)
			books = books[:h.CFG.OPDS.PAGE_SIZE-1]
		}
		if int(bc) > h.CFG.OPDS.PAGE_SIZE {
			if page > 1 {
				firstRef := fmt.Sprintf("/opds/search?language=%s&book=%s&page=1", lang, queryString)
				firstLink := &Link{Rel: FeedFirstLinkRel, Href: firstRef, Type: FeedNavigationLinkType}
				f.Link = append(f.Link, *firstLink)

				prevRef := fmt.Sprintf("/opds/search?language=%s&book=%s&page=%d", lang, queryString, page-1)
				prevLink := &Link{Rel: FeedPrevLinkRel, Href: prevRef, Type: FeedNavigationLinkType}
				f.Link = append(f.Link, *prevLink)
			}
			lastPage := int(math.Ceil(float64(bc) / float64(h.CFG.OPDS.PAGE_SIZE)))
			if page < lastPage {
				lastRef := fmt.Sprintf("/opds/search?language=%s&book=%s&page=%d", lang, queryString, lastPage)
				lastLink := &Link{Rel: FeedLastLinkRel, Href: lastRef, Type: FeedNavigationLinkType}
				f.Link = append(f.Link, *lastLink)
			}
		}

		h.feedBookEntries(r, books, f)
		writeFeed(w, http.StatusOK, *f)
	case ac != 0 && bc == 0: // show authors
		// h.listAuthors(w, r)
		page, err := strconv.Atoi(r.FormValue("page"))
		if err != nil {
			page = 1
		}
		offset := (page - 1) * h.CFG.OPDS.PAGE_SIZE
		authors := h.DB.PageFoundAuthors(queryString, h.CFG.OPDS.PAGE_SIZE+1, offset)
		selfHref = fmt.Sprintf("/opds/search?language=%s&author=%s&page=%d", lang, queryString, page)
		f := NewFeed(h.MP[lang].Sprintf("Found authors - %d", ac), "", selfHref)
		if len(authors) > h.CFG.OPDS.PAGE_SIZE {
			nextRef := fmt.Sprintf("/opds/search?language=%s&author=%s&page=%d", lang, queryString, page+1)
			nextLink := &Link{Rel: FeedNextLinkRel, Href: nextRef, Type: FeedNavigationLinkType}
			f.Link = append(f.Link, *nextLink)
			authors = authors[:h.CFG.OPDS.PAGE_SIZE-1]
		}
		if int(ac) > h.CFG.OPDS.PAGE_SIZE {
			if page > 1 {
				firstRef := fmt.Sprintf("/opds/search?language=%s&book=%s&page=1", lang, queryString)
				firstLink := &Link{Rel: FeedFirstLinkRel, Href: firstRef, Type: FeedNavigationLinkType}
				f.Link = append(f.Link, *firstLink)

				prevRef := fmt.Sprintf("/opds/search?language=%s&book=%s&page=%d", lang, queryString, page-1)
				prevLink := &Link{Rel: FeedPrevLinkRel, Href: prevRef, Type: FeedNavigationLinkType}
				f.Link = append(f.Link, *prevLink)
			}
			lastPage := int(math.Ceil(float64(ac) / float64(h.CFG.OPDS.PAGE_SIZE)))
			if page < lastPage {
				lastRef := fmt.Sprintf("/opds/search?language=%s&book=%s&page=%d", lang, queryString, lastPage)
				lastLink := &Link{Rel: FeedLastLinkRel, Href: lastRef, Type: FeedNavigationLinkType}
				f.Link = append(f.Link, *lastLink)
			}
		}

		// h.feedAuthorEntries(authors, f)
		for _, author := range authors {
			entry := &Entry{
				Title:   author.Sort,
				ID:      fmt.Sprintf("/opds/authors/author=%d", author.ID),
				Updated: f.Time(time.Now()),
				Links: []Link{
					{Rel: FeedSubsectionLinkRel, Href: fmt.Sprintf("/opds/authors?language=%s&id=%d", lang, author.ID), Type: FeedNavigationLinkType},
				},
				Content: &Content{
					Type:    FeedTextContentType,
					Content: h.MP[lang].Sprintf("Total books - %d", author.Count),
				},
			}
			f.Entry = append(f.Entry, entry)
		}
		writeFeed(w, http.StatusOK, *f)
	default:
		selfHref = fmt.Sprintf("/opds/search?language=%s&q={searchTerms}", lang)
		f := NewFeed(h.MP[lang].Sprintf("Nothing found"), "", selfHref)
		writeFeed(w, http.StatusOK, *f)
	}
}

// authors
func (h *Handler) authors(w http.ResponseWriter, r *http.Request) {
	switch {
	default: // Select author
		h.listAuthors(w, r)
		h.LOG.D.Println("ListAuthors")
	case r.FormValue("id") != "" && r.FormValue("anthology") == "" && r.FormValue("serie") == "": // Choose authors book select option: alphabetically or by series
		h.authorAnthology(w, r)
		h.LOG.D.Println("AuthorAnthology")
	case r.FormValue("id") != "" && r.FormValue("anthology") == "series": // Choose author book serie
		h.authorAnthologySeries(w, r)
		h.LOG.D.Println("AuthorAnthologySeries")
	case r.FormValue("id") != "" && (r.FormValue("anthology") == "alphabet" || r.FormValue("serie") != ""): // List all author books alphabetically or one serie books
		h.authorBooks(w, r)
		h.LOG.D.Println("AuthorBooks")
	}
}

// GET /opds/authors?author="" - all first authors letters
func (h *Handler) listAuthors(w http.ResponseWriter, r *http.Request) {
	lang := h.getLanguage(r)
	prefix := r.FormValue("author")
	abc := h.CFG.Locales.Languages[lang].Abc
	if r.Form.Has("all") {
		abc = ""
	}
	authors := h.DB.ListAuthors(prefix, abc)
	if len(authors) == 0 {
		return
	}
	sortAuthors(authors, h.CFG.Locales.Languages[lang].Tag)
	totalAuthors := 0
	for _, a := range authors {
		totalAuthors += a.Count
	}
	var selfHref string
	if prefix == "" {
		selfHref = fmt.Sprintf("/opds/authors?language=%s", lang)
		if abc == "" {
			selfHref += "&all"
		}
	} else {
		selfHref = fmt.Sprintf("/opds/authors?language=%s&author=%s", lang, url.QueryEscape(prefix))
	}

	f := NewFeed(h.MP[lang].Sprintf("Authors"), "", selfHref)
	addNotSpecLink := func() {
		if abc != "" {
			return
		}
		notSpecId := h.DB.AuthorNotSpecifiedId()
		if notSpecId > 0 {
			entry := &Entry{
				Title:   h.MP[lang].Sprintf("Author not specified"),
				ID:      fmt.Sprintf("/opds/authors/language=%s/id=%d", lang, notSpecId),
				Updated: f.Time(time.Now()),
				Links: []Link{
					{Rel: FeedSubsectionLinkRel, Href: fmt.Sprintf("/opds/authors?language=%s&id=%d", lang, notSpecId), Type: FeedNavigationLinkType},
				},
				Content: &Content{
					Type:    FeedTextContentType,
					Content: h.MP[lang].Sprintf("Author not specified"),
				},
			}
			f.Entry = append(f.Entry, entry)
		}
	}
	addAllAuthorsLinks := func() {
		if abc == "" || prefix != "" {
			return
		}
		entry := &Entry{
			Title:   h.MP[lang].Sprintf("All authors"),
			ID:      fmt.Sprintf("/opds/authors/language=%s/all", lang),
			Updated: f.Time(time.Now()),
			Links: []Link{
				{Rel: FeedSubsectionLinkRel, Href: fmt.Sprintf("/opds/authors?language=%s&all", lang), Type: FeedNavigationLinkType},
			},
			Content: &Content{
				Type:    FeedTextContentType,
				Content: h.MP[lang].Sprintf("Selection from all authors"),
			},
		}
		f.Entry = append(f.Entry, entry)
	}
	switch {
	case totalAuthors <= h.CFG.OPDS.PAGE_SIZE:
		authors = h.DB.ListAuthorWithTotals(prefix)
		for _, author := range authors {
			entry := &Entry{
				Title:   author.Sort,
				ID:      fmt.Sprintf("/opds/authors/language=%s/author=%d", lang, author.ID),
				Updated: f.Time(time.Now()),
				Links: []Link{
					{Rel: FeedSubsectionLinkRel, Href: fmt.Sprintf("/opds/authors?language=%s&id=%d", lang, author.ID), Type: FeedNavigationLinkType},
				},
				Content: &Content{
					Type:    FeedTextContentType,
					Content: h.MP[lang].Sprintf("Total books - %d", author.Count),
				},
			}
			f.Entry = append(f.Entry, entry)
		}
		// addAllAuthorsLinks()
		writeFeed(w, http.StatusOK, *f)
	default:
		for _, author := range authors {
			entry := &Entry{
				Title:   author.Sort,
				ID:      fmt.Sprintf("/opds/authors/language=%s/author=%s", lang, author.Sort),
				Updated: f.Time(time.Now()),
				Links: []Link{
					{Rel: FeedSubsectionLinkRel, Href: fmt.Sprintf("/opds/authors?language=%s&author=%s", lang, url.QueryEscape(author.Sort)), Type: FeedNavigationLinkType},
				},
				Content: &Content{
					Type:    FeedTextContentType,
					Content: h.MP[lang].Sprintf("Found authors - %d", author.Count),
				},
			}
			f.Entry = append(f.Entry, entry)
		}
		addNotSpecLink()
		addAllAuthorsLinks()
		writeFeed(w, http.StatusOK, *f)
	}
}

func (h *Handler) authorAnthology(w http.ResponseWriter, r *http.Request) {
	lang := h.getLanguage(r)
	authorId, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	authorSeries := h.DB.AuthorBookSeries(authorId)
	if len(authorSeries) > 0 {
		selfHref := fmt.Sprintf("/opds/authors?language=%s&id=%d", lang, authorId)
		author := h.fixIfNoSpecAuthorName(h.DB.AuthorByID(authorId), lang)
		f := NewFeed(author.Name, "", selfHref)
		f.Entry = []*Entry{
			{
				Title:   h.MP[lang].Sprintf("Alphabet"),
				ID:      fmt.Sprintf("/opds/authors/language=%s/id=%d/anthology=alphabet", lang, authorId),
				Updated: f.Time(time.Now()),
				Links: []Link{
					{
						Rel:  FeedSubsectionLinkRel,
						Href: fmt.Sprintf("/opds/authors?language=%s&id=%d&anthology=alphabet", lang, authorId),
						Type: FeedNavigationLinkType,
					},
				},
				Content: &Content{
					Type:    FeedTextContentType,
					Content: h.MP[lang].Sprintf("List books alphabetically"),
				},
			},
			{
				Title:   h.MP[lang].Sprintf("Series"),
				ID:      fmt.Sprintf("/opds/authors/language=%s/id=%d/anthology=series", lang, authorId),
				Updated: f.Time(time.Now()),
				Links: []Link{
					{
						Rel:  FeedSubsectionLinkRel,
						Href: fmt.Sprintf("/opds/authors?language=%s&id=%d&anthology=series", lang, authorId), Type: FeedNavigationLinkType,
					},
				},
				Content: &Content{
					Type:    FeedTextContentType,
					Content: h.MP[lang].Sprintf("List books series"),
				},
			},
		}
		writeFeed(w, http.StatusOK, *f)
	} else { // Author doesn't have book series
		h.authorBooks(w, r)
	}
}

func (h *Handler) authorAnthologySeries(w http.ResponseWriter, r *http.Request) {
	lang := h.getLanguage(r)
	authorId, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	author := h.fixIfNoSpecAuthorName(h.DB.AuthorByID(authorId), lang)
	selfHref := fmt.Sprintf("/opds/authors?language=%s&id=%d&anthology=series", lang, authorId)
	f := NewFeed(author.Name, "", selfHref)
	f.Entry = []*Entry{}
	var entry *Entry
	series := h.DB.AuthorBookSeries(authorId)
	for _, serie := range series {
		entry = &Entry{
			Title:   serie.Name,
			ID:      fmt.Sprintf("/opds/authors/language=%s/id=%d/serie=%d", lang, authorId, serie.ID),
			Updated: f.Time(time.Now()),
			Links: []Link{
				{Rel: FeedSubsectionLinkRel, Href: fmt.Sprintf("/opds/authors?language=%s&id=%d&serie=%d", lang, authorId, serie.ID), Type: FeedNavigationLinkType},
			},
		}
		f.Entry = append(f.Entry, entry)
	}
	writeFeed(w, http.StatusOK, *f)
}

func (h *Handler) authorBooks(w http.ResponseWriter, r *http.Request) {
	lang := h.getLanguage(r)
	authorId, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	serieId, _ := strconv.ParseInt(r.FormValue("serie"), 10, 64)
	author := h.fixIfNoSpecAuthorName(h.DB.AuthorByID(authorId), lang)
	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		page = 1
	}
	offset := (page - 1) * h.CFG.OPDS.PAGE_SIZE

	books := h.DB.ListAuthorBooks(authorId, serieId, h.CFG.OPDS.PAGE_SIZE+1, offset)
	selfHref := fmt.Sprintf("/opds/authors?language=%s&id=%d&anthology=alphabet&page=%d", lang, authorId, page)
	f := NewFeed(author.Name, "", selfHref)
	if len(books) > h.CFG.OPDS.PAGE_SIZE {
		nextRef := fmt.Sprintf("/opds/authors?language=%s&id=%d&anthology=alphabet&page=%d", lang, authorId, page+1)
		nextLink := &Link{Rel: FeedNextLinkRel, Href: nextRef, Type: FeedNavigationLinkType}
		f.Link = append(f.Link, *nextLink)
		books = books[:h.CFG.OPDS.PAGE_SIZE-1]
	}

	h.feedBookEntries(r, books, f)
	writeFeed(w, http.StatusOK, *f)
}

func (h *Handler) fixIfNoSpecAuthorName(author *model.Author, lang string) *model.Author {
	if author.Sort == "[author not specified]" || author.Name == "[author not specified]" {
		author.Name = h.MP[lang].Sprintf("Author not specified")
		author.Sort = h.MP[lang].Sprintf("Author not specified")
	}
	return author
}

// genres
func (h *Handler) genres(w http.ResponseWriter, r *http.Request) {
	switch {
	default:
		h.listGenres(w, r)
		h.LOG.D.Println("ListGenres")
	case r.FormValue("bunch") != "":
		h.listSubgenres(w, r)
		h.LOG.D.Println("ListSubgenres")
	case r.FormValue("code") != "":
		h.genreBooks(w, r)
		h.LOG.D.Println("GenreBooks")
	}
}

func (h *Handler) listGenres(w http.ResponseWriter, r *http.Request) {
	lang := h.getLanguage(r)
	selfHref := fmt.Sprintf("/opds/genres?language=%s", lang)
	f := NewFeed(h.MP[lang].Sprintf("Genres"), "", selfHref)
	f.Entry = []*Entry{}
	var entry *Entry
	genres := h.GT.ListGenres()
	title := ""
	content := ""
	for _, genre := range genres {
		for _, gd := range genre.Descriptions {
			if gd.Lang == lang {
				title = gd.Title
				content = gd.Detailed
				break
			}
		}
		if title != "" {
			entry = &Entry{
				Title:   title,
				ID:      fmt.Sprintf("/opds/genres/language=%s/bunch=%s", lang, genre.Value),
				Updated: f.Time(time.Now()),
				Links: []Link{
					{Rel: FeedSubsectionLinkRel, Href: fmt.Sprintf("/opds/genres?language=%s&bunch=%s", lang, genre.Value), Type: FeedNavigationLinkType},
				},
				Content: &Content{
					Content: content,
					Type:    FeedTextContentType,
				},
			}
		}
		f.Entry = append(f.Entry, entry)
	}
	writeFeed(w, http.StatusOK, *f)
}

func (h *Handler) listSubgenres(w http.ResponseWriter, r *http.Request) {
	lang := h.getLanguage(r)
	bunch := r.FormValue("bunch")
	selfHref := fmt.Sprintf("/opds/genres?language=%s&bunch=%s", lang, bunch)
	f := NewFeed(h.MP[lang].Sprintf("Genres"), "", selfHref)
	f.Entry = []*Entry{}
	var entry *Entry
	subgenres := h.GT.ListSubGenres(bunch)
	for _, sg := range subgenres {
		title := h.GT.SubgenreName(&sg, lang)
		gbc := h.DB.CountGenreBooks(sg.Value)
		if title != "" {
			entry = &Entry{
				Title:   title,
				ID:      fmt.Sprintf("/opds/genres/language=%s/code=%s", lang, sg.Value),
				Updated: f.Time(time.Now()),
				Links: []Link{
					{Rel: FeedSubsectionLinkRel, Href: fmt.Sprintf("/opds/genres?language=%s&code=%s", lang, sg.Value), Type: FeedAcquisitionLinkType},
				},
				Content: &Content{
					Content: h.MP[lang].Sprintf("Found titles - %d", gbc),
					Type:    FeedTextContentType,
				},
			}
		}
		f.Entry = append(f.Entry, entry)
	}
	writeFeed(w, http.StatusOK, *f)
}

func (h *Handler) genreBooks(w http.ResponseWriter, r *http.Request) {
	lang := h.getLanguage(r)
	genreCode := r.FormValue("code")
	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		page = 1
	}
	offset := (page - 1) * h.CFG.OPDS.PAGE_SIZE
	books := h.DB.ListGenreBooks(genreCode, h.CFG.OPDS.PAGE_SIZE+1, offset)
	selfHref := fmt.Sprintf("/opds/genres?language=%s&code=%s&page=%d", lang, genreCode, page)
	f := NewFeed(h.GT.GenreName(genreCode, h.getLanguage(r)), "", selfHref)
	if len(books) > h.CFG.OPDS.PAGE_SIZE {
		nextRef := fmt.Sprintf("/opds/genres?language=%s&code=%s&page=%d", lang, genreCode, page+1)
		nextLink := &Link{Rel: FeedNextLinkRel, Href: nextRef, Type: FeedNavigationLinkType}
		f.Link = append(f.Link, *nextLink)
		books = books[:h.CFG.OPDS.PAGE_SIZE]
	}
	if gbc := h.DB.CountGenreBooks(genreCode); int(gbc) > h.CFG.OPDS.PAGE_SIZE {
		if page > 1 {
			firstRef := fmt.Sprintf("/opds/genres?language=%s&code=%s&page=1", lang, genreCode)
			firstLink := &Link{Rel: FeedFirstLinkRel, Href: firstRef, Type: FeedNavigationLinkType}
			f.Link = append(f.Link, *firstLink)

			prevRef := fmt.Sprintf("/opds/genres?language=%s&code=%s&page=%d", lang, genreCode, page-1)
			prevLink := &Link{Rel: FeedPrevLinkRel, Href: prevRef, Type: FeedNavigationLinkType}
			f.Link = append(f.Link, *prevLink)
		}
		lastPage := int(math.Ceil(float64(gbc) / float64(h.CFG.OPDS.PAGE_SIZE)))
		if page < lastPage {
			lastRef := fmt.Sprintf("/opds/genres?language=%s&code=%s&page=%d", lang, genreCode, lastPage)
			lastLink := &Link{Rel: FeedLastLinkRel, Href: lastRef, Type: FeedNavigationLinkType}
			f.Link = append(f.Link, *lastLink)
		}
	}

	h.feedBookEntries(r, books, f)
	writeFeed(w, http.StatusOK, *f)
}

// series
func (h *Handler) series(w http.ResponseWriter, r *http.Request) {
	switch {
	default:
		h.listSeries(w, r)
		h.LOG.D.Println("listSeries")
	case r.FormValue("id") != "":
		h.serieBooks(w, r)
		h.LOG.D.Println("serieBooks")
	}
}

func (h *Handler) listSeries(w http.ResponseWriter, r *http.Request) {
	prefix := r.FormValue("serie")
	lang := h.getLanguage(r)
	series := h.DB.ListSeries(prefix, lang, h.CFG.Locales.Languages[lang].Abc)
	sortSeries(series, h.CFG.Locales.Languages[lang].Tag)
	if len(series) == 0 {
		return
	}
	totalBooks := 0
	for _, s := range series {
		totalBooks += s.Count
	}

	selfHref := ""
	if prefix == "" {
		selfHref = fmt.Sprintf("/opds/series?language=%s", lang)
	} else {
		selfHref = fmt.Sprintf("/opds/series?language=%s&serie=%s", lang, url.QueryEscape(prefix))
	}

	f := NewFeed(h.MP[lang].Sprintf("Series"), "", selfHref)
	switch {
	case len(series) <= h.CFG.OPDS.PAGE_SIZE && prefix != "":
		series = h.DB.ListSeriesWithTotals(prefix, lang)
		for _, serie := range series {
			entry := &Entry{
				Title:   serie.Name,
				ID:      fmt.Sprintf("/opds/series/language=%s/serie=%s", lang, serie.Name),
				Updated: f.Time(time.Now()),
				Links: []Link{
					{Rel: FeedSubsectionLinkRel, Href: fmt.Sprintf("/opds/series?language=%s&id=%d", lang, serie.ID), Type: FeedNavigationLinkType},
				},
				Content: &Content{
					Type:    FeedTextContentType,
					Content: h.MP[lang].Sprintf("Total books - %d", serie.Count),
				},
			}
			f.Entry = append(f.Entry, entry)
		}
		writeFeed(w, http.StatusOK, *f)
	default:
		for _, serie := range series {
			entry := &Entry{
				Title:   serie.Name,
				ID:      fmt.Sprintf("/opds/series/language=%s/serie=%s", lang, serie.Name),
				Updated: f.Time(time.Now()),
				Links: []Link{
					{Rel: FeedSubsectionLinkRel, Href: fmt.Sprintf("/opds/series?language=%s&serie=%s", lang, url.QueryEscape(serie.Name)), Type: FeedNavigationLinkType},
				},
				Content: &Content{
					Type:    FeedTextContentType,
					Content: h.MP[lang].Sprintf("Total series - %d", serie.Count),
				},
			}
			f.Entry = append(f.Entry, entry)
		}
		writeFeed(w, http.StatusOK, *f)
	}
}

func (h *Handler) serieBooks(w http.ResponseWriter, r *http.Request) {
	lang := h.getLanguage(r)
	serieId, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	serie := h.DB.SerieByID(serieId)
	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		page = 1
	}
	offset := (page - 1) * h.CFG.OPDS.PAGE_SIZE

	books := h.DB.ListSerieBooks(serieId, h.CFG.OPDS.PAGE_SIZE+1, offset)
	selfHref := fmt.Sprintf("/opds/series?language=%s&id=%d&page=%d", lang, serieId, page)
	f := NewFeed(serie.Name, "", selfHref)
	if len(books) > h.CFG.OPDS.PAGE_SIZE {
		nextRef := fmt.Sprintf("/opds/series?language=%s&id=%d&page=%d", lang, serieId, page+1)
		nextLink := &Link{Rel: FeedNextLinkRel, Href: nextRef, Type: FeedNavigationLinkType}
		f.Link = append(f.Link, *nextLink)
		books = books[:h.CFG.OPDS.PAGE_SIZE-1]
	}

	h.feedBookEntries(r, books, f)
	writeFeed(w, http.StatusOK, *f)
}

// Books
func (h *Handler) books(w http.ResponseWriter, r *http.Request) {
	switch {
	default:
	case r.FormValue("id") != "":
		h.unloadBook(w, r)
		h.LOG.D.Println("UnloadBook")
	}
}

func (h *Handler) feedBookEntries(r *http.Request, books []*model.Book, f *Feed) {
	lang := h.getLanguage(r)
	for _, book := range books {
		var authorsList []Author
		var authorsLinks []Link
		authors := h.DB.AuthorsByBookId(book.ID)
		for _, a := range authors {
			a = h.fixIfNoSpecAuthorName(a, lang)
			author := Author{
				Name: a.Name,
			}
			authorLink := Link{
				Title: fmt.Sprintf("%s - %s", h.MP[lang].Sprintf("All author books"), a.Name),
				Rel:   FeedRelatedLinkRel,
				Href:  fmt.Sprintf("/opds/authors?language=%s&id=%d", lang, a.ID),
				Type:  FeedNavigationLinkType,
			}

			authorsList = append(authorsList, author)
			authorsLinks = append(authorsLinks, authorLink)
		}

		links := append(authorsLinks, h.acquisitionLinks(book)...)
		entry := &Entry{
			Title:   book.Title,
			ID:      fmt.Sprintf("/opds/books/id=%d", book.ID),
			Updated: f.Time(time.Now()),
			Links:   links,
			Authors: authorsList,
			Content: &Content{
				Type:    FeedHtmlContentType,
				Content: h.contentInfo(r, book),
			},
		}
		f.Entry = append(f.Entry, entry)
	}
}

func (h *Handler) acquisitionLinks(book *model.Book) []Link {
	rel := "http://opds-spec.org/acquisition/open-access"
	link := []Link{}
	switch book.Format {
	case "fb2":
		linkFunc := func(convert string) Link {
			return Link{
				Rel:  rel,
				Href: fmt.Sprintf("/opds/books?id=%d&convert=%s", book.ID, convert),
				Type: mime.TypeByExtension(fmt.Sprintf(".fb2.%s", convert)),
			}
		}
		if h.CFG.OPDS.NO_CONVERSION {
			link = append(link, linkFunc("zip"))
		} else {
			link = append(link, linkFunc("epub"), linkFunc("zip"))
		}
	default:
		link = append(link,
			Link{
				Rel:  rel,
				Href: fmt.Sprintf("/opds/books?id=%d", book.ID),
				Type: mime.TypeByExtension("." + book.Format),
			},
		)

	}
	if book.Cover != "" {
		link = append(link,
			Link{
				Rel:  "http://opds-spec.org/image",
				Href: fmt.Sprintf("/opds/covers?cover=%d", book.ID),
				Type: mime.TypeByExtension(path.Ext(book.Cover)),
			},
		)
		link = append(link,
			Link{
				Rel:  "http://opds-spec.org/image/thumbnail",
				Href: fmt.Sprintf("/opds/covers?thumbnail=%d", book.ID),
				Type: mime.TypeByExtension(path.Ext(book.Cover)),
			},
		)
	}
	return link
}

func (h *Handler) unloadBook(w http.ResponseWriter, r *http.Request) {
	lang := h.getLanguage(r)
	bookId, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	book := h.DB.FindBookById(bookId)
	if book == nil {
		writeMessage(w, http.StatusNotFound, h.MP[lang].Sprintf("Book not found"))
		return
	}
	convert := r.FormValue("convert")
	ext := ""
	switch convert {
	case "epub":
		ext = ".epub"
	case "zip":
		ext = ".zip"
	}

	var rc io.ReadCloser
	if book.Archive == "" {
		rc, _ = os.Open(path.Join(h.CFG.Library.STOCK_DIR, book.File))
	} else {
		zr, _ := zip.OpenReader(path.Join(h.CFG.Library.STOCK_DIR, book.Archive))
		defer zr.Close()
		for _, file := range zr.File {
			if file.Name == book.File {
				rc, _ = file.Open()
				break
			}
		}
	}
	defer rc.Close()

	// w.Header().Add("Content-Type", fmt.Sprintf("%s; name=%s", mime.TypeByExtension("." + book.Format + zipExt), book.File+zipExt))
	w.Header().Add("Content-Type", mime.TypeByExtension("."+book.Format+ext))
	w.Header().Add("Content-Transfer-Encoding", "binary")
	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", book.File+ext))
	w.WriteHeader(http.StatusOK)

	switch convert {
	case "epub":
		rsc, err := NewReadSeekCloser(rc)
		if err != nil {
			h.LOG.E.Println(err)
			return
		}
		wc := NewWriteCloser(w)
		err = h.ConvertFb2Epub(bookId, rsc, wc)
		if err != nil {
			h.LOG.E.Println(err)
		}
	case "zip":
		zipWriter := zip.NewWriter(w)
		defer zipWriter.Close()
		fileWriter, _ := zipWriter.CreateHeader(
			&zip.FileHeader{
				Name:   book.File,
				Method: zip.Deflate,
			},
		)
		io.Copy(fileWriter, rc)
		zipWriter.Flush()
	default:
		io.Copy(w, rc)
	}
}

// Covers
func (h *Handler) covers(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.FormValue("cover") != "":
		h.LOG.D.Println(commentURL("Cover", r))
		h.unloadCover(w, r)
	case r.FormValue("thumbnail") != "":
		h.LOG.D.Println(commentURL("Thumbnail", r))
		h.unloadThumbnail(w, r)
	default:
		return
	}

}

func (h *Handler) unloadCover(w http.ResponseWriter, r *http.Request) {
	bookId, _ := strconv.ParseInt(r.FormValue("cover"), 10, 64)
	img := h.getCoverImage(bookId)
	if img == nil {
		return
	}
	w.Header().Add("Content-Disposition", "attachment; filename=cover.jpg")
	w.Header().Add("Content-Type", "image/jpeg")
	jpeg.Encode(w, img, nil)
}

func (h *Handler) unloadThumbnail(w http.ResponseWriter, r *http.Request) {
	bookId, _ := strconv.ParseInt(r.FormValue("thumbnail"), 10, 64)
	img := h.getCoverImage(bookId)
	if img == nil {
		return
	}
	img = resize.Resize(100, 0, img, resize.NearestNeighbor)
	// img = imaging.Resize(img, 100, 0, imaging.NearestNeighbor)
	w.Header().Add("Content-Disposition", "attachment; filename=thumbnail.jpg")
	w.Header().Add("Content-Type", "image/jpeg")
	jpeg.Encode(w, img, nil)
}

func (h *Handler) getCoverImage(bookId int64) (img image.Image) {
	book := h.DB.FindBookById(bookId)
	if book == nil {
		return nil
	}
	if book.Cover == "" {
		return nil
	}

	switch book.Format {
	case "fb2":
		img, err := fb2.GetCoverImage(h.CFG.Library.STOCK_DIR, book)
		if err != nil {
			h.LOG.D.Print(err)
			return nil
		}
		return img
	case "epub":
		img, err := epub.GetCoverImage(h.CFG.Library.STOCK_DIR, book)
		if err != nil {
			h.LOG.D.Print(err)
			return nil
		}
		return img
	}
	return nil
}

func sortAuthors(s []*model.Author, t language.Tag) {
	c := collate.New(t, collate.Force)
	sort.Slice(s, func(i, j int) bool {
		return c.CompareString(s[i].Sort, s[j].Sort) < 0
	})
}

func sortSeries(s []*model.Serie, t language.Tag) {
	c := collate.New(t, collate.Force)
	sort.Slice(s, func(i, j int) bool {
		return c.CompareString(s[i].Name, s[j].Name) < 0
	})
}

func (h *Handler) contentInfo(r *http.Request, b *model.Book) (info string) {
	lang := h.getLanguage(r)
	if b.Plot != "" {
		info += fmt.Sprintf("<p>%s</p>", b.Plot)
	}
	if b.Year != "0" {
		info += fmt.Sprintf("%s: %s<br/>", h.MP[lang].Sprintf("Year"), b.Year)
	}
	info += fmt.Sprintf("%s: %d Kb<br/>", h.MP[lang].Sprintf("Size"), int(float32(b.Size)/1024))
	if b.Serie.Name != "" {
		info += fmt.Sprintf("%s: %s", h.MP[lang].Sprintf("Serie"), b.Serie.Name)
		if b.SerieNum > 0 {
			info += fmt.Sprintf(" #%d", b.SerieNum)
		}
		info += "<br/>"
	}
	return
}

func (h *Handler) getLanguage(r *http.Request) string {
	lang := r.FormValue("language")
	if lang == "" {
		lang = h.CFG.Locales.DEFAULT
	}
	if strings.Contains(h.CFG.Locales.ACCEPTED, lang) {
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

func (h *Handler) ConvertFb2Epub(b int64, r io.ReadSeekCloser, w io.WriteCloser) error {
	fb := &cfb2.FB2Parser{
		BookId:  b,
		LOG:     h.LOG,
		DB:      h.DB,
		RC:      r,
		Decoder: u8xml.NewDecoder(r),
	}

	if err := fb.MakeEpub(w); err != nil {
		return err
	}
	return nil
}

type ResponseWriteCloser struct {
	http.ResponseWriter
}

func NewWriteCloser(w http.ResponseWriter) *ResponseWriteCloser {
	return &ResponseWriteCloser{
		ResponseWriter: w,
	}
}

func (w ResponseWriteCloser) Write(b []byte) (int, error) {
	return w.ResponseWriter.Write(b)
}

func (w ResponseWriteCloser) Close() error {
	return nil
}

type BufferedReadSeekCloser struct {
	io.ReadSeeker
}

func NewReadSeekCloser(r io.ReadCloser) (*BufferedReadSeekCloser, error) {
	if rs, ok := r.(io.ReadSeeker); ok {
		return &BufferedReadSeekCloser{
			ReadSeeker: rs,
		}, nil
	}
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	rs := bytes.NewReader(b)

	return &BufferedReadSeekCloser{
		ReadSeeker: rs,
	}, nil
}

func (r BufferedReadSeekCloser) Close() error {
	return nil
}
