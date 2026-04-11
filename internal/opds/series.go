package opds

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"

	"github.com/vinser/flibgolite/internal/core/model"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

// Series
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
	var (
		abc   string
		aLang string
		all   string
	)
	if r.Form.Has("all") {
		abc = ""
		aLang = ""
		all = "&all"
	} else {
		abc = h.CFG.Languages[lang].Abc + `'0','1','2','3','4','5','6','7','8','9','0'`
		aLang = lang
	}
	series := h.DB.ListSeries(prefix, aLang, abc)
	if len(series) == 0 {
		return
	}
	sortSeries(series, h.CFG.Locales.Languages[lang].Tag)
	totalSeries := 0
	for _, s := range series {
		totalSeries += s.Count
	}

	selfHref := ""
	if prefix == "" {
		selfHref = fmt.Sprintf("/opds/series?language=%s", lang)
	} else {
		selfHref = fmt.Sprintf("/opds/series?language=%s%s&serie=%s", lang, all, url.QueryEscape(prefix))
	}

	f := NewFeed(h.MP[lang].Sprintf("Series"), "", selfHref)
	addAllSeriesLink := func() {
		if abc == "" || prefix != "" {
			return
		}
		entry := &Entry{
			Title:   h.MP[lang].Sprintf("~All series"),
			ID:      fmt.Sprintf("/opds/series/language=%s/all", lang),
			Updated: f.Time(time.Now()),
			Links: []Link{
				{Rel: FeedSubsectionLinkRel, Href: fmt.Sprintf("/opds/series?language=%s&all", lang), Type: FeedNavigationLinkType},
			},
			Content: &Content{
				Type:    FeedTextContentType,
				Content: h.MP[lang].Sprintf("^Selection from all series"),
			},
		}
		f.Entry = append(f.Entry, entry)
	}
	switch {
	case totalSeries <= h.CFG.OPDS.PAGE_SIZE:
		series = h.DB.ListSeriesWithTotals(prefix, aLang)
		for _, serie := range series {
			entry := &Entry{
				Title:   serie.Name,
				ID:      fmt.Sprintf("/opds/series/language=%s/serie=%s", lang, serie.Name),
				Updated: f.Time(time.Now()),
				Links: []Link{
					{Rel: FeedSubsectionLinkRel, Href: fmt.Sprintf("/opds/series?language=%s%s&id=%d", lang, all, serie.ID), Type: FeedNavigationLinkType},
				},
				Content: &Content{
					Type:    FeedTextContentType,
					Content: h.MP[lang].Sprintf("^Series Total books - %d", serie.Count),
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
					{Rel: FeedSubsectionLinkRel, Href: fmt.Sprintf("/opds/series?language=%s%s&serie=%s", lang, all, url.QueryEscape(serie.Name)), Type: FeedNavigationLinkType},
				},
				Content: &Content{
					Type:    FeedTextContentType,
					Content: h.MP[lang].Sprintf("^Total series - %d", serie.Count),
				},
			}
			f.Entry = append(f.Entry, entry)
		}
		addAllSeriesLink()
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

func sortSeries(s []*model.Serie, t language.Tag) {
	c := collate.New(t, collate.Force)
	sort.Slice(s, func(i, j int) bool {
		return c.CompareString(s[i].Name, s[j].Name) < 0
	})
}
