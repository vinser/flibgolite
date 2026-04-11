package opds

import (
	"fmt"
	"net/http"
	"sort"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language/display"
)

// Languages
func (h *Handler) languages(w http.ResponseWriter, r *http.Request) {
	lang := h.getLanguage(r)
	selfHref := fmt.Sprintf("/opds/languages?language=%s", lang)
	f := NewFeed(h.MP[lang].Sprintf("Choose language"), "", selfHref)
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
				Content: h.MP[lang].Sprintf("^Language Total books - %d", langBookTotal),
			},
		}
		f.Entry = append(f.Entry, entry)
	}
	writeFeed(w, http.StatusOK, *f)
}
