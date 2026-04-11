package opds

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/vinser/flibgolite/internal/core/model"

	_ "image/gif"
	_ "image/png"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

// Authors
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
	abc := h.CFG.Languages[lang].Abc
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
		if utf8.RuneCountInString(prefix) > 0 {
			return
		}
		notSpecId := h.DB.AuthorNotSpecifiedId()
		if notSpecId > 0 {
			entry := &Entry{
				Title:   h.MP[lang].Sprintf("~Author not specified"),
				ID:      fmt.Sprintf("/opds/authors/language=%s/id=%d", lang, notSpecId),
				Updated: f.Time(time.Now()),
				Links: []Link{
					{Rel: FeedSubsectionLinkRel, Href: fmt.Sprintf("/opds/authors?language=%s&id=%d", lang, notSpecId), Type: FeedNavigationLinkType},
				},
				Content: &Content{
					Type:    FeedTextContentType,
					Content: h.MP[lang].Sprintf("~Author not specified"),
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
			Title:   h.MP[lang].Sprintf("~All authors"),
			ID:      fmt.Sprintf("/opds/authors/language=%s/all", lang),
			Updated: f.Time(time.Now()),
			Links: []Link{
				{Rel: FeedSubsectionLinkRel, Href: fmt.Sprintf("/opds/authors?language=%s&all", lang), Type: FeedNavigationLinkType},
			},
			Content: &Content{
				Type:    FeedTextContentType,
				Content: h.MP[lang].Sprintf("^Selection from all authors"),
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
					Content: h.MP[lang].Sprintf("^Author Total books - %d", author.Count),
				},
			}
			f.Entry = append(f.Entry, entry)
		}
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
					Content: h.MP[lang].Sprintf("^Found authors - %d", author.Count),
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
				Title:   h.MP[lang].Sprintf("~Alphabet"),
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
					Content: h.MP[lang].Sprintf("^List books alphabetically"),
				},
			},
			{
				Title:   h.MP[lang].Sprintf("~Author Series"),
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
					Content: h.MP[lang].Sprintf("^List books series"),
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
		author.Name = h.MP[lang].Sprintf("~Author not specified")
		author.Sort = h.MP[lang].Sprintf("~Author not specified")
	}
	return author
}

func sortAuthors(s []*model.Author, t language.Tag) {
	c := collate.New(t, collate.Force)
	sort.Slice(s, func(i, j int) bool {
		return c.CompareString(s[i].Sort, s[j].Sort) < 0
	})
}
