package opds

import (
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"time"
	"unicode/utf8"
)

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

func (h *Handler) foundAuthorsEntry(f *Feed, lang, queryString string, authorCount int64) *Entry {
	return &Entry{
		Title:   h.MP[lang].Sprintf("~Authors found"),
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
			Content: h.MP[lang].Sprintf("^Authors found - %d", authorCount),
		},
	}
}

func (h *Handler) foundBooksEntry(f *Feed, lang, queryString string, titleCount int64) *Entry {
	return &Entry{
		Title:   h.MP[lang].Sprintf("~Titles"),
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
			Content: h.MP[lang].Sprintf("^Found books - %d", titleCount),
		},
	}
}

func (h *Handler) foundKeywordsEntry(f *Feed, lang, queryString string, keywordCount int64) *Entry {
	return &Entry{
		Title:   h.MP[lang].Sprintf("~Keywords"),
		ID:      fmt.Sprintf("/opds/search/keywords=%s", queryString),
		Updated: f.Time(time.Now()),
		Links: []Link{
			{
				Rel:  FeedSubsectionLinkRel,
				Href: fmt.Sprintf("/opds/search?language=%s&keywords=%s", lang, queryString),
				Type: FeedNavigationLinkType,
			},
		},
		Content: &Content{
			Type:    FeedTextContentType,
			Content: h.MP[lang].Sprintf("^Found by keywords - %d", keywordCount),
		},
	}
}
func (h *Handler) serach(w http.ResponseWriter, r *http.Request) {
	lang := h.getLanguage(r)
	h.LOG.D.Println(commentURL("Search", r))
	selfHref := ""
	queryString := ""
	var authorCount, titleCount, keywordCount int64
	switch {
	case r.FormValue("q") != "":
		queryString = r.FormValue("q")
		if utf8.RuneCountInString(queryString) < 3 {
			authorCount = 0
			titleCount = 0
			keywordCount = 0
		}
		authorCount = h.DB.SearchAuthorsCount(queryString)
		titleCount = h.DB.SearchBooksCountByTitle(queryString)
		keywordCount = h.DB.SearchBooksCountByKeyword(queryString)
	case r.FormValue("author") != "":
		queryString = r.FormValue("author")
		authorCount = h.DB.SearchAuthorsCount(queryString)
	case r.FormValue("book") != "":
		queryString = r.FormValue("book")
		titleCount = h.DB.SearchBooksCountByTitle(queryString)
	case r.FormValue("keywords") != "":
		queryString = r.FormValue("keywords")
		keywordCount = h.DB.SearchBooksCountByKeyword(queryString)
	}
	switch {
	case (authorCount == 0 && titleCount == 0 && keywordCount == 0): // nothing found
		selfHref = fmt.Sprintf("/opds/search?language=%s&q={searchTerms}", lang)
		f := NewFeed(h.MP[lang].Sprintf("Nothing found"), "", selfHref)
		writeFeed(w, http.StatusOK, *f)
	case authorCount > 0 && titleCount == 0 && keywordCount == 0: // show found authors
		// h.listAuthors(w, r)
		page, err := strconv.Atoi(r.FormValue("page"))
		if err != nil {
			page = 1
		}
		offset := (page - 1) * h.CFG.OPDS.PAGE_SIZE
		authors := h.DB.PageFoundAuthors(queryString, h.CFG.OPDS.PAGE_SIZE+1, offset)
		selfHref = fmt.Sprintf("/opds/search?language=%s&author=%s&page=%d", lang, queryString, page)
		f := NewFeed(h.MP[lang].Sprintf("Found authors - %d", authorCount), "", selfHref)
		if len(authors) > h.CFG.OPDS.PAGE_SIZE {
			nextRef := fmt.Sprintf("/opds/search?language=%s&author=%s&page=%d", lang, queryString, page+1)
			nextLink := &Link{Rel: FeedNextLinkRel, Href: nextRef, Type: FeedNavigationLinkType}
			f.Link = append(f.Link, *nextLink)
			authors = authors[:h.CFG.OPDS.PAGE_SIZE-1]
		}
		if int(authorCount) > h.CFG.OPDS.PAGE_SIZE {
			if page > 1 {
				firstRef := fmt.Sprintf("/opds/search?language=%s&book=%s&page=1", lang, queryString)
				firstLink := &Link{Rel: FeedFirstLinkRel, Href: firstRef, Type: FeedNavigationLinkType}
				f.Link = append(f.Link, *firstLink)

				prevRef := fmt.Sprintf("/opds/search?language=%s&book=%s&page=%d", lang, queryString, page-1)
				prevLink := &Link{Rel: FeedPrevLinkRel, Href: prevRef, Type: FeedNavigationLinkType}
				f.Link = append(f.Link, *prevLink)
			}
			lastPage := int(math.Ceil(float64(authorCount) / float64(h.CFG.OPDS.PAGE_SIZE)))
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
					Content: h.MP[lang].Sprintf("^Total books found - %d", author.Count),
				},
			}
			f.Entry = append(f.Entry, entry)
		}
		writeFeed(w, http.StatusOK, *f)
	case authorCount == 0 && titleCount > 0 && keywordCount == 0: // show books found by title
		page, err := strconv.Atoi(r.FormValue("page"))
		if err != nil {
			page = 1
		}
		offset := (page - 1) * h.CFG.OPDS.PAGE_SIZE
		books := h.DB.PageFoundBooksByTitle(queryString, h.CFG.OPDS.PAGE_SIZE+1, offset)
		selfHref = fmt.Sprintf("/opds/search?language=%s&book=%s&page=%d", lang, queryString, page)
		f := NewFeed(h.MP[lang].Sprintf("Found books - %d", titleCount), "", selfHref)
		if len(books) > h.CFG.OPDS.PAGE_SIZE {
			nextRef := fmt.Sprintf("/opds/search?language=%s&book=%s&page=%d", lang, queryString, page+1)
			nextLink := &Link{Rel: FeedNextLinkRel, Href: nextRef, Type: FeedNavigationLinkType}
			f.Link = append(f.Link, *nextLink)
			books = books[:h.CFG.OPDS.PAGE_SIZE-1]
		}
		if int(titleCount) > h.CFG.OPDS.PAGE_SIZE {
			if page > 1 {
				firstRef := fmt.Sprintf("/opds/search?language=%s&book=%s&page=1", lang, queryString)
				firstLink := &Link{Rel: FeedFirstLinkRel, Href: firstRef, Type: FeedNavigationLinkType}
				f.Link = append(f.Link, *firstLink)

				prevRef := fmt.Sprintf("/opds/search?language=%s&book=%s&page=%d", lang, queryString, page-1)
				prevLink := &Link{Rel: FeedPrevLinkRel, Href: prevRef, Type: FeedNavigationLinkType}
				f.Link = append(f.Link, *prevLink)
			}
			lastPage := int(math.Ceil(float64(titleCount) / float64(h.CFG.OPDS.PAGE_SIZE)))
			if page < lastPage {
				lastRef := fmt.Sprintf("/opds/search?language=%s&book=%s&page=%d", lang, queryString, lastPage)
				lastLink := &Link{Rel: FeedLastLinkRel, Href: lastRef, Type: FeedNavigationLinkType}
				f.Link = append(f.Link, *lastLink)
			}
		}

		h.feedBookEntries(r, books, f)
		writeFeed(w, http.StatusOK, *f)
	case authorCount == 0 && titleCount == 0 && keywordCount > 0: // show books found by keyword
		page, err := strconv.Atoi(r.FormValue("page"))
		if err != nil {
			page = 1
		}
		offset := (page - 1) * h.CFG.OPDS.PAGE_SIZE
		books := h.DB.PageFoundBooksByKeywords(queryString, h.CFG.OPDS.PAGE_SIZE+1, offset)
		selfHref = fmt.Sprintf("/opds/search?language=%s&keywords=%s&page=%d", lang, queryString, page)
		f := NewFeed(h.MP[lang].Sprintf("Found by keywords - %d", keywordCount), "", selfHref)
		if len(books) > h.CFG.OPDS.PAGE_SIZE {
			nextRef := fmt.Sprintf("/opds/search?language=%s&keywords=%s&page=%d", lang, queryString, page+1)
			nextLink := &Link{Rel: FeedNextLinkRel, Href: nextRef, Type: FeedNavigationLinkType}
			f.Link = append(f.Link, *nextLink)
			books = books[:h.CFG.OPDS.PAGE_SIZE-1]
		}
		if int(keywordCount) > h.CFG.OPDS.PAGE_SIZE {
			if page > 1 {
				firstRef := fmt.Sprintf("/opds/search?language=%s&keywords=%s&page=1", lang, queryString)
				firstLink := &Link{Rel: FeedFirstLinkRel, Href: firstRef, Type: FeedNavigationLinkType}
				f.Link = append(f.Link, *firstLink)

				prevRef := fmt.Sprintf("/opds/search?language=%s&keywords=%s&page=%d", lang, queryString, page-1)
				prevLink := &Link{Rel: FeedPrevLinkRel, Href: prevRef, Type: FeedNavigationLinkType}
				f.Link = append(f.Link, *prevLink)
			}
			lastPage := int(math.Ceil(float64(keywordCount) / float64(h.CFG.OPDS.PAGE_SIZE)))
			if page < lastPage {
				lastRef := fmt.Sprintf("/opds/search?language=%s&keywords=%s&page=%d", lang, queryString, lastPage)
				lastLink := &Link{Rel: FeedLastLinkRel, Href: lastRef, Type: FeedNavigationLinkType}
				f.Link = append(f.Link, *lastLink)
			}
		}

		h.feedBookEntries(r, books, f)
		writeFeed(w, http.StatusOK, *f)
	default: // show chices for found items
		selfHref = fmt.Sprintf("/opds/search?language=%s&q={searchTerms}", lang)
		f := NewFeed(h.MP[lang].Sprintf("Choose from the found ones"), "", selfHref)
		f.Entry = []*Entry{}
		if authorCount > 0 {
			f.Entry = append(f.Entry, h.foundAuthorsEntry(f, lang, queryString, authorCount))
		}
		if titleCount > 0 {
			f.Entry = append(f.Entry, h.foundBooksEntry(f, lang, queryString, titleCount))
		}
		if keywordCount > 0 {
			f.Entry = append(f.Entry, h.foundKeywordsEntry(f, lang, queryString, keywordCount))
		}
		writeFeed(w, http.StatusOK, *f)
	}
}
