package opds

import (
	"bytes"
	_ "embed"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

const (
	// Feed link types
	FeedAcquisitionLinkType       = "application/atom+xml;profile=opds-catalog;kind=acquisition"
	FeedNavigationLinkType        = "application/atom+xml;profile=opds-catalog;kind=navigation"
	FeedSearchDescriptionLinkType = "application/opensearchdescription+xml"

	// Feed link relations
	FeedStartLinkRel      = "start"
	FeedSelfLinkRel       = "self"
	FeedSearchLinkRel     = "search"
	FeedFirstLinkRel      = "first"
	FeedLastLinkRel       = "last"
	FeedNextLinkRel       = "next"
	FeedPrevLinkRel       = "previous"
	FeedSubsectionLinkRel = "subsection"
	FeedRelatedLinkRel    = "related"

	// Content types
	FeedTextContentType     = "text"
	FeedHtmlContentType     = "html"
	FeedTextHtmlContentType = "text/html"
)

type Feed struct {
	XMLName      xml.Name `xml:"feed"`
	Xmlns        string   `xml:"xmlns,attr"`
	XmlnsDC      string   `xml:"xmlns:dc,attr,omitempty"`
	XmlnsOS      string   `xml:"xmlns:os,attr,omitempty"`
	XmlnsOPDS    string   `xml:"xmlns:opds,attr,omitempty"`
	Title        string   `xml:"title"`
	ID           string   `xml:"id"`
	Updated      TimeStr  `xml:"updated"`
	Icon         string   `xml:"icon,omitempty"`
	Link         []Link   `xml:"link"`
	Author       []Author `xml:"author,omitempty"`
	Entry        []*Entry `xml:"entry"`
	Category     string   `xml:"category,omitempty"`
	Logo         string   `xml:"logo,omitempty"`
	Content      string   `xml:"content,omitempty"`
	Subtitle     string   `xml:"subtitle,omitempty"`
	SearchResult uint     `xml:"opensearch:totalResults,omitempty"`
}

type Entry struct {
	// XMLName   xml.Name `xml:"entry"`
	// Xmlns     string   `xml:"xmlns,attr,omitempty"`
	Title     string   `xml:"title"`
	ID        string   `xml:"id"`
	Links     []Link   `xml:"link"`
	Published string   `xml:"published,omitempty"`
	Updated   TimeStr  `xml:"updated"`
	Category  string   `xml:"category,omitempty"`
	Authors   []Author `xml:"author"`
	Summary   *Summary `xml:"summary"`
	Content   *Content `xml:"content"`
	Rights    string   `xml:"rights,omitempty"`
	Source    string   `xml:"source,omitempty"`
}

type Link struct {
	// XMLName xml.Name `xml:"link"`
	Type   string `xml:"type,attr,omitempty"`
	Title  string `xml:"title,attr,omitempty"`
	Href   string `xml:"href,attr"`
	Rel    string `xml:"rel,attr,omitempty"`
	Length string `xml:"length,attr,omitempty"`
}

type Author struct {
	// XMLName xml.Name `xml:"author"`
	Name string `xml:"name,omitempty"`
	Uri  string `xml:"uri,omitempty"`
}

type Summary struct {
	// XMLName xml.Name `xml:"summary"`
	Content string `xml:",chardata"`
	Type    string `xml:"type,attr"`
}

type Content struct {
	// XMLName xml.Name `xml:"content"`
	Content string `xml:",chardata"`
	Type    string `xml:"type,attr"`
}

type TimeStr string

func (f *Feed) Time(t time.Time) TimeStr {
	return TimeStr(t.Format("2006-01-02T15:04:05-07:00"))
}

var idReplace = regexp.MustCompile(`\&amp;|\?|\&`)

//go:embed favicon.ico
var FAVICON_ICO []byte

func NewFeed(title, subtitle, self string) *Feed {
	f := &Feed{
		XMLName:   xml.Name{},
		Xmlns:     "http://www.w3.org/2005/Atom",
		XmlnsDC:   "http://purl.org/dc/terms/",
		XmlnsOS:   "http://a9.com/-/spec/opensearch/1.1/",
		XmlnsOPDS: "http://opds-spec.org/2010/catalog",
		Title:     title,
		Icon:      "/favicon.ico",
		ID:        idReplace.ReplaceAllString(self, "/"),
		Link: []Link{
			{Rel: FeedStartLinkRel, Href: "/opds", Type: FeedNavigationLinkType},
			{Rel: FeedSelfLinkRel, Href: self, Type: FeedNavigationLinkType},
		},
		Subtitle: subtitle,
	}
	f.Updated = f.Time(time.Now())
	return f
}

func commentURL(comment string, r *http.Request) string {
	qu, _ := url.QueryUnescape(r.URL.String())
	return fmt.Sprintf("%s --->URL: [%s]", comment, qu)
}

func writeFeed(w http.ResponseWriter, statusCode int, f Feed) {
	data, err := xml.MarshalIndent(f, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "Internal server error")
		return
	}
	s := fmt.Sprintf("%s%s", xml.Header, data)
	w.Header().Add("Content-Type", "application/atom+xml;charset=utf-8")
	w.WriteHeader(statusCode)
	io.WriteString(w, s)
}

func writeMessage(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	io.WriteString(w, message)
}

// favicon.ico
func (h *Handler) unloadFavicon(w http.ResponseWriter) {
	w.Header().Add("Content-Disposition", "attachment; filename=favicon.ico")
	w.Header().Add("Content-Type", "image/x-icon")
	w.WriteHeader(http.StatusOK)
	io.Copy(w, bytes.NewReader(FAVICON_ICO))
}
