package opds

import (
	"archive/zip"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"math"
	"mime"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/nfnt/resize"
	cfb2 "github.com/vinser/flibgolite/internal/converter/fb2"
	"github.com/vinser/flibgolite/internal/core/model"
	"github.com/vinser/flibgolite/internal/parsers"
	"github.com/vinser/flibgolite/internal/parsers/epub"
	"github.com/vinser/flibgolite/internal/parsers/fb2"
	"github.com/vinser/u8xml"

	"github.com/mozillazg/go-unidecode"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

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
		f := NewFeed(h.MP[lang].Sprintf("Latest Found titles - %d", bc), "", selfHref)
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
		f := NewFeed(h.MP[lang].Sprintf("Nothing new"), "", selfHref)
		writeFeed(w, http.StatusOK, *f)
	}
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
				Title: fmt.Sprintf("%s - %s", h.MP[lang].Sprintf("~All author books"), a.Name),
				Rel:   FeedRelatedLinkRel,
				Href:  fmt.Sprintf("/opds/authors?language=%s&id=%d", lang, a.ID),
				Type:  FeedNavigationLinkType,
			}

			authorsList = append(authorsList, author)
			authorsLinks = append(authorsLinks, authorLink)
		}

		links := append(authorsLinks, h.acquisitionLinks(book)...)
		if serie := h.DB.SerieByBookID(book.ID); serie != nil {
			serieLink := Link{
				Title: fmt.Sprintf("%s - %s", h.MP[lang].Sprintf("~All serie books"), serie.Name),
				Rel:   FeedRelatedLinkRel,
				Href:  fmt.Sprintf("/opds/series?language=%s&id=%d&page=1", lang, serie.ID),
				Type:  FeedNavigationLinkType,
			}
			links = append(links, serieLink)
		}

		bookLang := ""
		if book.Language != nil && book.Language.Code != "" {
			bookLang = book.Language.Code
		}

		bookYear := ""
		if book.Year != "" && book.Year != "0" {
			bookYear = book.Year
		}

		entry := &Entry{
			Title:      book.Title,
			ID:         fmt.Sprintf("/opds/books/id=%d", book.ID),
			Updated:    f.Time(time.Now()),
			Links:      links,
			Authors:    authorsList,
			DcLanguage: bookLang,
			DcIssued:   bookYear,
			Content: &Content{
				Type:    FeedTextHtmlContentType,
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

	convert := r.FormValue("convert")
	ext := ""
	switch convert {
	case "epub":
		ext = ".epub"
	case "zip":
		ext = ".zip"
	}
	authors := h.DB.AuthorsByBookId(bookId)
	authorName := ""
	switch {
	case len(authors) == 0:
		authorName = "Author not specified"
	case len(authors) > 1:
		authorName = "Group of authors"
	default:
		authorName = authors[0].Sort
	}

	// w.Header().Add("Content-Type", fmt.Sprintf("%s; name=%s", mime.TypeByExtension("." + book.Format + zipExt), book.File+zipExt))
	w.Header().Add("Content-Type", mime.TypeByExtension("."+book.Format+ext))
	w.Header().Add("Content-Transfer-Encoding", "binary")
	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileNameByAuthorTitle(authorName, book.Title)+"."+book.Format+ext))
	w.WriteHeader(http.StatusOK)

	switch convert {
	case "epub":
		rsc, err := NewReadSeekCloser(rc)
		if err != nil {
			h.LOG.E.Println(err)
			return
		}
		wc := NewWriteCloser(w)
		err = h.ConvertFb2Epub(wc, rsc, bookId)
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

// Transliterated file name
// RegExp Find illegal file name characters
var rxNotFileName = regexp.MustCompile(`[^0-9a-zA-Z-_]`)

const MAX_FILE_NAME_LEN = 232

func fileNameByAuthorTitle(author, title string) string {
	if author != "" {
		names := strings.Split(unidecode.Unidecode(parsers.CollapseSpaces(strings.ReplaceAll(author, ",", " "))), " ")
		for i := range names {
			if names[i] != "" {
				names[i] = strings.ToLower(names[i])
				names[i] = strings.ToUpper(names[i][:1]) + names[i][1:]
			}
		}
		author = strings.Join(names, "-")
	}
	if title != "" {
		words := strings.Split(unidecode.Unidecode(strings.ReplaceAll(parsers.CollapseSpaces(title), ",", " ")), " ")
		title = strings.Join(words, "-")
	}
	fileName := rxNotFileName.ReplaceAllString(unidecode.Unidecode(parsers.CollapseSpaces(author+"_"+title)), "")
	switch {
	case len(fileName) == 0:
		fileName = "book"
	case len(fileName) > MAX_FILE_NAME_LEN:
		fileName = fileName[:MAX_FILE_NAME_LEN]
	}
	return fileName
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

func (h *Handler) ConvertFb2Epub(w io.WriteCloser, r io.ReadSeekCloser, b int64) error {
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

// Info
func (h *Handler) contentInfo(r *http.Request, b *model.Book) (info string) {
	lang := h.getLanguage(r)
	info = "<div>"
	if b.Plot != "" {
		info += fmt.Sprintf("<p>%s</p>", b.Plot)
	}
	if b.Language.Code != "" {
		info += fmt.Sprintf("<br/>%s: %s", h.MP[lang].Sprintf("Language"), cases.Title(language.Make(b.Language.Code)).String(display.Self.Name(language.Make(b.Language.Code))))
	}
	if b.Year != "0" {
		info += fmt.Sprintf("<br/>%s: %s", h.MP[lang].Sprintf("Year"), b.Year)
	}
	if b.Archive != "" {
		info += fmt.Sprintf("<br/>%s: %s", h.MP[lang].Sprintf("Archive"), b.Archive)
	}
	info += fmt.Sprintf("<br/>%s: %s", h.MP[lang].Sprintf("File"), b.File)
	info += fmt.Sprintf("<br/>%s: %d Kb", h.MP[lang].Sprintf("Size"), int(float32(b.Size)/1024))
	if b.Serie.Name != "" {
		info += fmt.Sprintf("<br/>%s: %s", h.MP[lang].Sprintf("Serie"), b.Serie.Name)
		if b.SerieNum > 0 {
			info += fmt.Sprintf(" #%d", b.SerieNum)
		}
		info += "<br/>"
	}
	return info + "</div>"
}
