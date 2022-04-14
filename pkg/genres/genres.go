package genres

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type GenresTree struct {
	XMLName xml.Name `xml:"fbgenrestransfer"`
	Genres  []Genre  `xml:"genre"`
}

type Genre struct {
	Value        string      `xml:"value,attr"`
	Descriptions []RootDescr `xml:"root-descr"`
	Subgenres    []Subgenre  `xml:"subgenres>subgenre"`
}

type RootDescr struct {
	Lang     string `xml:"lang,attr"`
	Title    string `xml:"genre-title,attr"`
	Detailed string `xml:"detailed,attr"`
}

type Subgenre struct {
	Value        string       `xml:"value,attr"`
	Descriptions []GenreDescr `xml:"genre-descr"`
	Alts         []GenreAlt   `xml:"genre-alt"`
}

type GenreDescr struct {
	Lang  string `xml:"lang,attr"`
	Title string `xml:"title,attr"`
}

type GenreAlt struct {
	Value  string `xml:"value,attr"`
	Format string `xml:"format,attr"`
}

func NewGenresTree(treeFile string) *GenresTree {
	var err error
	gt := &GenresTree{}
	xmlStream, err := os.Open(treeFile)
	if err != nil {
		// log.Println("failed to open genres tree file")
		return gt
	}
	defer xmlStream.Close()
	xmlData, _ := ioutil.ReadAll(xmlStream)
	decoder := xml.NewDecoder(bytes.NewReader(xmlData))
	decoder.Strict = false
	decoder.Decode(&gt)
	return gt
}

func (gt *GenresTree) Transfer(genre string) string {
	genre = strings.ReplaceAll(strings.TrimSpace(genre), "-", "_")
	for _, g := range gt.Genres {
		for _, sg := range g.Subgenres {
			if sg.Value == genre {
				return sg.Value
			}
			for _, sga := range sg.Alts {
				if sga.Value == genre {
					return sg.Value
				}
			}
		}
	}
	return string([]rune("no_name:" + genre)[:32])
	// return genre
}

func (gt *GenresTree) ListGenres() []Genre {
	return gt.Genres
}

func (gt *GenresTree) ListSubGenres(genre string) []Subgenre {
	sg := []Subgenre{}
	for _, g := range gt.Genres {
		if g.Value == genre {
			return g.Subgenres
		}
	}
	return sg
}

func (gt *GenresTree) GenreName(genre, lang string) string {
	for _, g := range gt.Genres {
		for _, sg := range g.Subgenres {
			if sg.Value == genre {
				for _, sgd := range sg.Descriptions {
					if sgd.Lang == lang {
						return sgd.Title
					}
				}
			}
		}
	}
	return ""
}

func (gt *GenresTree) SubgenreName(sg *Subgenre, lang string) string {
	for _, sgd := range sg.Descriptions {
		if sgd.Lang == lang {
			return sgd.Title
		}
	}
	return ""
}

func (gt *GenresTree) String() string {
	s := ""
	for _, g := range gt.Genres {
		s += fmt.Sprintln("genre:", g.Value)
		for _, rd := range g.Descriptions {
			s += fmt.Sprintf("\troot-descr: lang=%s title=%s detailed=%s\n", rd.Lang, rd.Title, rd.Detailed)
		}
		s += fmt.Sprintln("\tsubgenres ---")
		for _, sg := range g.Subgenres {
			s += fmt.Sprintln("\t\tsubgenre: value=", sg.Value)
			for _, sgd := range sg.Descriptions {
				s += fmt.Sprintf("\t\t\tgenre-descr lang=%s title=%s\n", sgd.Lang, sgd.Title)
			}
			for _, sga := range sg.Alts {
				s += fmt.Sprintf("\t\t\tgenre-alt value=%s format=%s\n", sga.Value, sga.Format)
			}
		}
	}
	return s
}
