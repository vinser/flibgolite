package opds

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"
)

// Genres
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
					Content: h.MP[lang].Sprintf("^Genres Found titles - %d", gbc),
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
	books := h.DB.PageGenreBooks(genreCode, h.CFG.OPDS.PAGE_SIZE+1, offset)
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
