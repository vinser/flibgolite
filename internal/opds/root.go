package opds

import (
	"fmt"
	"net/http"
	"time"
)

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
			Title:   h.MP[lang].Sprintf("~Latest Books"),
			ID:      "latest",
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
				Content: h.MP[lang].Sprintf("^Browse the latest books received"),
			},
		},
		{
			Title:   h.MP[lang].Sprintf("~Book Authors"),
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
				Content: h.MP[lang].Sprintf("^Browse books by author"),
			},
		},
		{
			Title:   h.MP[lang].Sprintf("~Book Series"),
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
				Content: h.MP[lang].Sprintf("^Browse books by series"),
			},
		},
		{
			Title:   h.MP[lang].Sprintf("~Book Genres"),
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
				Content: h.MP[lang].Sprintf("^Browse books by genre"),
			},
		},
	}
	if len(h.CFG.Languages) > 1 {
		f.Entry = append(f.Entry, &Entry{
			Title:   h.MP[lang].Sprintf("~Book Languages"),
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
				Content: h.MP[lang].Sprintf("^Language selection"),
			},
		})
	}

	writeFeed(w, http.StatusOK, *f)
}
