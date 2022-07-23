package opds

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"image"
	"image/jpeg"
	"io"
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

	"github.com/vinser/flibgolite/pkg/config"
	"github.com/vinser/flibgolite/pkg/database"
	"github.com/vinser/flibgolite/pkg/epub"
	"github.com/vinser/flibgolite/pkg/fb2"
	"github.com/vinser/flibgolite/pkg/genres"
	"github.com/vinser/flibgolite/pkg/model"
	"github.com/vinser/flibgolite/pkg/rlog"

	"github.com/nfnt/resize"
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
	P   *message.Printer
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.LOG.I.Println(commentURL("Router", r))
	// switch r.URL.Path {
	switch strings.ReplaceAll(r.URL.Path, "//", "/") { // compensate PocketBook Reader search query error
	case "/opds":
		h.root(w, r)
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
	selfHref := "/opds"
	lang := r.FormValue("language")
	if lang == "" {
		lang = h.CFG.Locales.DEFAULT
	}
	f := NewFeed(h.CFG.OPDS.TITLE, "", selfHref)
	f.Entry = []*Entry{
		{
			Title:   h.P.Sprintf("Book Authors"),
			ID:      "authors",
			Updated: f.Time(time.Now()),
			Link: []Link{
				{Rel: FeedSubsectionLinkRel, Href: "/opds/authors?language=" + lang, Type: FeedNavigationLinkType},
			},
			Content: &Content{
				Type:    FeedTextContentType,
				Content: h.P.Sprintf("Choose an author of a book"),
			},
		},
		{
			Title:   h.P.Sprintf("Book Genres"),
			ID:      "genres",
			Updated: f.Time(time.Now()),
			Link: []Link{
				{Rel: FeedSubsectionLinkRel, Href: "/opds/genres?language=" + lang, Type: FeedNavigationLinkType},
			},
			Content: &Content{
				Type:    FeedTextContentType,
				Content: h.P.Sprintf("Choose a genre of a book"),
			},
		},
		{
			Title:   h.P.Sprintf("Book Series"),
			ID:      "series",
			Updated: f.Time(time.Now()),
			Link: []Link{
				{Rel: FeedSubsectionLinkRel, Href: "/opds/series?language=" + lang, Type: FeedNavigationLinkType},
			},
			Content: &Content{
				Type:    FeedTextContentType,
				Content: h.P.Sprintf("Choose a serie of a book"),
			},
		},
		{
			Title:   h.P.Sprintf("Book Languages"),
			ID:      "languages",
			Updated: f.Time(time.Now()),
			Link: []Link{
				{Rel: FeedSubsectionLinkRel, Href: "/opds/languages", Type: FeedNavigationLinkType},
			},
			Content: &Content{
				Type:    FeedTextContentType,
				Content: h.P.Sprintf("Choose book language"),
			},
		},
	}
	//
	writeFeed(w, http.StatusOK, *f)
}

// Languages
func (h *Handler) languages(w http.ResponseWriter, r *http.Request) {
	selfHref := "/opds/languages"
	f := NewFeed("Choose book language", "", selfHref)
	ordered := []string{}
	for o := range h.CFG.Locales.LANG {
		ordered = append(ordered, o)
	}
	sort.Strings(ordered)
	for _, k := range ordered {
		langTitle := cases.Title(h.CFG.Locales.LANG[k].Tag)
		langName := langTitle.String(display.Self.Name(h.CFG.Locales.LANG[k].Tag))
		entry := &Entry{
			Title:   langName,
			ID:      "/opds?language=" + k,
			Updated: f.Time(time.Now()),
			Link: []Link{
				{Rel: FeedSubsectionLinkRel, Href: "/opds?language=" + k, Type: FeedNavigationLinkType},
			},
		}
		f.Entry = append(f.Entry, entry)
	}
	writeFeed(w, http.StatusOK, *f)
}

// OpenSearch description document
func (h *Handler) openSerach(w http.ResponseWriter, r *http.Request) {
	data :=
		`
<OpenSearchDescription xmlns="http://a9.com/-/spec/opensearch/1.1/">
<ShortName>FLib Go Go Go!!!</ShortName>
<Description>Search on catalog</Description>
<InputEncoding>UTF-8</InputEncoding>
<OutputEncoding>UTF-8</OutputEncoding>
<Url type="application/atom+xml;profile=opds-catalog;kind=acquisition" template="/opds/search?q={searchTerms}"/>
</OpenSearchDescription>	
`
	s := fmt.Sprintf("%s%s", xml.Header, data)
	w.Header().Add("Content-Type", "application/atom+xml")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, s)

}

// Search
func (h *Handler) serach(w http.ResponseWriter, r *http.Request) {
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
		selfHref = "/opds/search?q={searchTerms}"
		f := NewFeed(h.P.Sprintf("Choose from the found ones"), "", selfHref)
		f.Entry = []*Entry{
			{
				Title:   h.P.Sprintf("Titles"),
				ID:      fmt.Sprint("/opds/search?book=", queryString),
				Updated: f.Time(time.Now()),
				Link: []Link{
					{Rel: FeedSubsectionLinkRel, Href: fmt.Sprint("/opds/search?book=", queryString), Type: FeedNavigationLinkType},
				},
				Content: &Content{
					Type:    FeedTextContentType,
					Content: h.P.Sprintf("Found titles - %d", bc),
				},
			},
			{
				Title:   h.P.Sprintf("Authors"),
				ID:      fmt.Sprint("/opds/search?author=", queryString),
				Updated: f.Time(time.Now()),
				Link: []Link{
					{Rel: FeedSubsectionLinkRel, Href: fmt.Sprint("/opds/search?author=", queryString), Type: FeedNavigationLinkType},
				},
				Content: &Content{
					Type:    FeedTextContentType,
					Content: h.P.Sprintf("Found authors - %d", ac),
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
		selfHref = fmt.Sprintf("/opds/search?book=%s&page=%d", queryString, page)
		f := NewFeed(h.P.Sprintf("Found titles - %d", bc), "", selfHref)
		if len(books) > h.CFG.OPDS.PAGE_SIZE {
			nextRef := fmt.Sprintf("/opds/search?book=%s&page=%d", queryString, page+1)
			nextLink := &Link{Rel: FeedNextLinkRel, Href: nextRef, Type: FeedNavigationLinkType}
			f.Link = append(f.Link, *nextLink)
			books = books[:h.CFG.OPDS.PAGE_SIZE-1]
		}
		h.feedBookEntries(books, f)
		writeFeed(w, http.StatusOK, *f)
	case ac != 0 && bc == 0: // show authors
		// h.listAuthors(w, r)
		page, err := strconv.Atoi(r.FormValue("page"))
		if err != nil {
			page = 1
		}
		offset := (page - 1) * h.CFG.OPDS.PAGE_SIZE
		authors := h.DB.PageFoundAuthors(queryString, h.CFG.OPDS.PAGE_SIZE+1, offset)
		selfHref = fmt.Sprintf("/opds/search?author=%s&page=%d", queryString, page)
		f := NewFeed(h.P.Sprintf("Found titles - %d", bc), "", selfHref)
		if len(authors) > h.CFG.OPDS.PAGE_SIZE {
			nextRef := fmt.Sprintf("/opds/search?author=%s&page=%d", queryString, page+1)
			nextLink := &Link{Rel: FeedNextLinkRel, Href: nextRef, Type: FeedNavigationLinkType}
			f.Link = append(f.Link, *nextLink)
			authors = authors[:h.CFG.OPDS.PAGE_SIZE-1]
		}
		// h.feedAuthorEntries(authors, f)
		for _, author := range authors {
			entry := &Entry{
				Title:   author.Sort,
				ID:      "/opds/authors?author=" + author.Sort,
				Updated: f.Time(time.Now()),
				Link: []Link{
					{Rel: FeedSubsectionLinkRel, Href: "/opds/authors?author=" + url.QueryEscape(author.Sort), Type: FeedNavigationLinkType},
				},
				Content: &Content{
					Type:    FeedTextContentType,
					Content: h.P.Sprintf("Found authors - %d", author.Count),
				},
			}
			f.Entry = append(f.Entry, entry)
		}
		writeFeed(w, http.StatusOK, *f)
	default:
		return
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
	prefix := r.FormValue("author")
	lang := r.FormValue("language")
	if lang == "" {
		lang = h.CFG.Locales.DEFAULT
	}
	authors := h.DB.ListAuthors(prefix, h.CFG.Locales.LANG[lang].Abc)
	sortAuthors(authors, h.CFG.Locales.LANG[lang].Tag)
	if len(authors) == 0 {
		return
	}
	totalAuthors := 0
	for _, a := range authors {
		totalAuthors += a.Count
	}

	var selfHref string
	if prefix == "" {
		selfHref = "/opds/authors"
	} else {
		selfHref = "/opds/authors?author=" + url.QueryEscape(prefix)
	}

	f := NewFeed(h.P.Sprintf("Authors"), "", selfHref)
	switch {
	case totalAuthors <= h.CFG.OPDS.PAGE_SIZE:
		authors = h.DB.ListAuthorWithTotals(prefix)
		for i := range authors {
			entry := &Entry{
				Title:   authors[i].Sort,
				ID:      "/opds/authors?author=" + authors[i].Sort,
				Updated: f.Time(time.Now()),
				Link: []Link{
					{Rel: FeedSubsectionLinkRel, Href: "/opds/authors?id=" + fmt.Sprint(authors[i].ID), Type: FeedNavigationLinkType},
				},
				Content: &Content{
					Type:    FeedTextContentType,
					Content: h.P.Sprintf("Total books - %d", authors[i].Count),
				},
			}
			f.Entry = append(f.Entry, entry)
		}
		writeFeed(w, http.StatusOK, *f)
	default:
		for i := range authors {
			entry := &Entry{
				Title:   authors[i].Sort,
				ID:      "/opds/authors?author=" + authors[i].Sort,
				Updated: f.Time(time.Now()),
				Link: []Link{
					{Rel: FeedSubsectionLinkRel, Href: "/opds/authors?author=" + url.QueryEscape(authors[i].Sort), Type: FeedNavigationLinkType},
				},
				Content: &Content{
					Type:    FeedTextContentType,
					Content: h.P.Sprintf("Found authors - %d", authors[i].Count),
				},
			}
			f.Entry = append(f.Entry, entry)
		}
		writeFeed(w, http.StatusOK, *f)
	}
}

// GET /opds/authors?id="" - all first authors letters
func (h *Handler) authorAnthology(w http.ResponseWriter, r *http.Request) {
	authorId, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	authorSeries := h.DB.AuthorBookSeries(authorId)
	if len(authorSeries) > 0 {
		selfHref := "/opds/authors?id=" + r.FormValue("id")
		author := h.DB.AuthorByID(authorId)
		f := NewFeed(author.Name, "", selfHref)
		f.Entry = []*Entry{
			{
				Title:   h.P.Sprintf("Alphabet"),
				ID:      fmt.Sprint("/opds/authors?id=", authorId, "&anthology=alphabet"),
				Updated: f.Time(time.Now()),
				Link: []Link{
					{Rel: FeedSubsectionLinkRel, Href: fmt.Sprint("/opds/authors?id=", authorId, "&anthology=alphabet"), Type: FeedNavigationLinkType},
				},
				Content: &Content{
					Type:    FeedTextContentType,
					Content: h.P.Sprintf("List books alphabetically"),
				},
			},
			{
				Title:   h.P.Sprintf("Series"),
				ID:      fmt.Sprint("/opds/authors?id=", authorId, "&anthology=series"),
				Updated: f.Time(time.Now()),
				Link: []Link{
					{Rel: FeedSubsectionLinkRel, Href: fmt.Sprint("/opds/authors?id=", authorId, "&anthology=series"), Type: FeedNavigationLinkType},
				},
				Content: &Content{
					Type:    FeedTextContentType,
					Content: h.P.Sprintf("List books series"),
				},
			},
		}
		writeFeed(w, http.StatusOK, *f)
	} else { // Author doesn't have book series
		h.authorBooks(w, r)
	}
}

func (h *Handler) authorAnthologySeries(w http.ResponseWriter, r *http.Request) {
	authorId, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	author := h.DB.AuthorByID(authorId)
	selfHref := fmt.Sprint("/opds/authors?id=", authorId, "&anthology=series")
	f := NewFeed(author.Name, "", selfHref)
	f.Entry = []*Entry{}
	var entry *Entry
	series := h.DB.AuthorBookSeries(authorId)
	for _, serie := range series {
		entry = &Entry{
			Title:   serie.Name,
			ID:      fmt.Sprint("/opds/authors?id=", authorId, "&serie=", serie.ID),
			Updated: f.Time(time.Now()),
			Link: []Link{
				{Rel: FeedSubsectionLinkRel, Href: fmt.Sprint("/opds/authors?id=", authorId, "&serie=", serie.ID), Type: FeedNavigationLinkType},
			},
		}
		f.Entry = append(f.Entry, entry)
	}
	writeFeed(w, http.StatusOK, *f)
}

func (h *Handler) authorBooks(w http.ResponseWriter, r *http.Request) {
	authorId, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	serieId, _ := strconv.ParseInt(r.FormValue("serie"), 10, 64)
	author := h.DB.AuthorByID(authorId)
	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		page = 1
	}
	offset := (page - 1) * h.CFG.OPDS.PAGE_SIZE

	books := h.DB.ListAuthorBooks(authorId, serieId, h.CFG.OPDS.PAGE_SIZE+1, offset)
	selfHref := fmt.Sprintf("/opds/authors?id=%d&anthology=alphabet&page=%d", authorId, page)
	f := NewFeed(author.Name, "", selfHref)
	if len(books) > h.CFG.OPDS.PAGE_SIZE {
		nextRef := fmt.Sprintf("/opds/authors?id=%d&anthology=alphabet&page=%d", authorId, page+1)
		nextLink := &Link{Rel: FeedNextLinkRel, Href: nextRef, Type: FeedNavigationLinkType}
		f.Link = append(f.Link, *nextLink)
		books = books[:h.CFG.OPDS.PAGE_SIZE-1]
	}

	h.feedBookEntries(books, f)
	writeFeed(w, http.StatusOK, *f)
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
	selfHref := "/opds/genres"
	f := NewFeed(h.P.Sprintf("Genres"), "", selfHref)
	f.Entry = []*Entry{}
	var entry *Entry
	genres := h.GT.ListGenres()
	title := ""
	content := ""
	for _, genre := range genres {
		for _, gd := range genre.Descriptions {
			if gd.Lang == h.CFG.Locales.DEFAULT {
				title = gd.Title
				content = gd.Detailed
				break
			}
		}
		if title != "" {
			entry = &Entry{
				Title:   title,
				ID:      fmt.Sprint("/opds/genres?bunch=", genre.Value),
				Updated: f.Time(time.Now()),
				Link: []Link{
					{Rel: FeedSubsectionLinkRel, Href: fmt.Sprint("/opds/genres?bunch=", genre.Value), Type: FeedNavigationLinkType},
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
	bunch := r.FormValue("bunch")
	selfHref := fmt.Sprint("/opds/genres?bunch=", bunch)
	f := NewFeed(h.P.Sprintf("Genres"), "", selfHref)
	f.Entry = []*Entry{}
	var entry *Entry
	subgenres := h.GT.ListSubGenres(bunch)
	for _, sg := range subgenres {
		title := h.GT.SubgenreName(&sg, h.CFG.Locales.DEFAULT)
		gbc := h.DB.CountGenreBooks(sg.Value)
		if title != "" {
			entry = &Entry{
				Title:   title,
				ID:      fmt.Sprint("/opds/genres?code=", sg.Value),
				Updated: f.Time(time.Now()),
				Link: []Link{
					{Rel: FeedSubsectionLinkRel, Href: fmt.Sprint("/opds/genres?code=", sg.Value), Type: FeedAcquisitionLinkType},
				},
				Content: &Content{
					Content: h.P.Sprintf("Found titles - %d", gbc),
					Type:    FeedTextContentType,
				},
			}
		}
		f.Entry = append(f.Entry, entry)
	}
	writeFeed(w, http.StatusOK, *f)
}

func (h *Handler) genreBooks(w http.ResponseWriter, r *http.Request) {
	genreCode := r.FormValue("code")
	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		page = 1
	}
	offset := (page - 1) * h.CFG.OPDS.PAGE_SIZE
	books := h.DB.ListGenreBooks(genreCode, h.CFG.OPDS.PAGE_SIZE+1, offset)
	selfHref := fmt.Sprintf("/opds/genres?code=%s&page=%d", genreCode, page)
	f := NewFeed(h.GT.GenreName(genreCode, h.CFG.Locales.DEFAULT), "", selfHref)
	if len(books) > h.CFG.OPDS.PAGE_SIZE {
		nextRef := fmt.Sprintf("/opds/genres?code=%s&page=%d", genreCode, page+1)
		nextLink := &Link{Rel: FeedNextLinkRel, Href: nextRef, Type: FeedNavigationLinkType}
		f.Link = append(f.Link, *nextLink)
		books = books[:h.CFG.OPDS.PAGE_SIZE-1]
	}

	h.feedBookEntries(books, f)
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
	lang := r.FormValue("language")
	if lang == "" {
		lang = h.CFG.Locales.DEFAULT
	}
	series := h.DB.ListSeries(prefix, lang, h.CFG.Locales.LANG[lang].Abc)
	sortSeries(series, h.CFG.Locales.LANG[lang].Tag)
	if len(series) == 0 {
		return
	}
	totalBooks := 0
	for _, s := range series {
		totalBooks += s.Count
	}

	selfHref := ""
	if prefix == "" {
		selfHref = "/opds/series?language=" + lang
	} else {
		selfHref = "/opds/series?serie=" + url.QueryEscape(prefix) + "&language=" + lang
	}

	f := NewFeed(h.P.Sprintf("Series"), "", selfHref)
	switch {
	case len(series) <= h.CFG.OPDS.PAGE_SIZE && prefix != "":
		series = h.DB.ListSeriesWithTotals(prefix, lang)
		for _, serie := range series {
			entry := &Entry{
				Title:   serie.Name,
				ID:      "/opds/series?serie=" + serie.Name,
				Updated: f.Time(time.Now()),
				Link: []Link{
					{Rel: FeedSubsectionLinkRel, Href: "/opds/series?id=" + fmt.Sprint(serie.ID), Type: FeedNavigationLinkType},
				},
				Content: &Content{
					Type:    FeedTextContentType,
					Content: h.P.Sprintf("Total books - %d", serie.Count),
				},
			}
			f.Entry = append(f.Entry, entry)
		}
		writeFeed(w, http.StatusOK, *f)
	default:
		for _, serie := range series {
			entry := &Entry{
				Title:   serie.Name,
				ID:      "/opds/series?serie=" + serie.Name + "&language=" + lang,
				Updated: f.Time(time.Now()),
				Link: []Link{
					{Rel: FeedSubsectionLinkRel, Href: "/opds/series?serie=" + url.QueryEscape(serie.Name) + "&language=" + lang, Type: FeedNavigationLinkType},
				},
				Content: &Content{
					Type:    FeedTextContentType,
					Content: h.P.Sprintf("Total series - %d", serie.Count),
				},
			}
			f.Entry = append(f.Entry, entry)
		}
		writeFeed(w, http.StatusOK, *f)
	}
}

func (h *Handler) serieBooks(w http.ResponseWriter, r *http.Request) {
	serieId, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	serie := h.DB.SerieByID(serieId)
	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		page = 1
	}
	offset := (page - 1) * h.CFG.OPDS.PAGE_SIZE

	books := h.DB.ListSerieBooks(serieId, h.CFG.OPDS.PAGE_SIZE+1, offset)
	selfHref := fmt.Sprintf("/opds/series?id=%d&page=%d", serieId, page)
	f := NewFeed(serie.Name, "", selfHref)
	if len(books) > h.CFG.OPDS.PAGE_SIZE {
		nextRef := fmt.Sprintf("/opds/series?id=%d&page=%d", serieId, page+1)
		nextLink := &Link{Rel: FeedNextLinkRel, Href: nextRef, Type: FeedNavigationLinkType}
		f.Link = append(f.Link, *nextLink)
		books = books[:h.CFG.OPDS.PAGE_SIZE-1]
	}

	h.feedBookEntries(books, f)
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

func (h *Handler) feedBookEntries(books []*model.Book, f *Feed) {
	for _, book := range books {
		authors := h.DB.AuthorsByBookId(book.ID)
		author := ""
		for _, a := range authors {
			author += fmt.Sprint(a.Name, ", ")
		}
		entry := &Entry{
			Title:   book.Title,
			ID:      fmt.Sprint("/opds/books?id=", book.ID),
			Updated: f.Time(time.Now()),
			Link: []Link{
				{
					Rel:  "http://opds-spec.org/acquisition/open-access",
					Href: fmt.Sprint("/opds/books?id=", book.ID),
					Type: fmt.Sprint("application/", book.Format),
				},
				{
					Rel:  "http://opds-spec.org/image",
					Href: fmt.Sprint("/opds/covers?cover=", book.ID),
					Type: mime.TypeByExtension(path.Ext(book.Cover)),
				},
				{
					Rel:  "http://opds-spec.org/image/thumbnail",
					Href: fmt.Sprint("/opds/covers?thumbnail=", book.ID),
					Type: mime.TypeByExtension(path.Ext(book.Cover)),
				},
			},
			Authors: []Author{
				{
					Name: strings.TrimSuffix(author, ", "),
				},
			},
			Content: &Content{
				Type:    FeedHtmlContentType,
				Content: fmt.Sprint(book.Plot),
			},
		}
		f.Entry = append(f.Entry, entry)
	}
}

func (h *Handler) unloadBook(w http.ResponseWriter, r *http.Request) {
	bookId, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	book := h.DB.FindBookById(bookId)
	if book == nil {
		writeMessage(w, http.StatusNotFound, h.P.Sprintf("Book not found"))
		return
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

	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", book.File))
	w.Header().Add("Content-Type", fmt.Sprintf("application/fb2; name=%s", book.File))
	w.Header().Add("Content-Transfer-Encoding", "binary")
	w.WriteHeader(http.StatusOK)
	io.Copy(w, rc)
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
	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "cover.jpg"))
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
	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "thumbnail.jpg"))
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
			h.LOG.E.Print(err)
			return nil
		}
		return img
	case "epub":
		img, err := epub.GetCoverImage(h.CFG.Library.STOCK_DIR, book)
		if err != nil {
			h.LOG.E.Print(err)
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
