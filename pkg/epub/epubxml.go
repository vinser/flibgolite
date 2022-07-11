package epub

import (
	"archive/zip"
	"bufio"
	"encoding/xml"
	"fmt"
	"image"
	"io"
	"log"
	"path"
	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/vinser/flibgolite/pkg/model"
	"golang.org/x/net/html/charset"
)

// Container ----------------------------------

// EPUB Open Container Format (OCF) 3.2
type OCF struct {
	Container xml.Name `xml:"container"`
	RootFiles struct {
		RootFile struct {
			FullPath string `xml:"full-path,attr"`
		} `xml:"rootfile"`
	} `xml:"rootfiles"`
}

// Get OPF file path fron OCF file
func GetOPFPath(zr *zip.ReadCloser) (string, error) {
	f, _ := zr.Open("META-INF/container.xml")
	ocf := &OCF{}
	if err := decodeXML(f, &ocf); err != nil {
		return "", err
	}
	return ocf.RootFiles.RootFile.FullPath, nil
}

// Packages ------------------------------------

// EPUB Packages 3.2
type OPF struct {
	// The root element of the Package Document and defines various aspects of the EPUB Package
	Package xml.Name `xml:"package"`
	// Specifies the language used in the contents and attribute values of the carrying element and its descendants
	Lang string `xml:"xml:lang,attr,omitempty"`
	// Specifies the EPUB specification version to which the given EPUB Package conforms.
	Version string `xml:"version,attr"`

	// Here only a minimal set of meta information for Reading Systems to use to internally catalogue an EPUB Publication
	Metadata struct {
		// Contains an identifier associated with the given Rendition, such as a UUID, DOI or ISBN.
		Identifier []struct {
			ID string `xml:"id,attr,omitempty"`
			// ID value
			Text string `xml:",chardata"`
		} `xml:"identifier"`
		// Rrepresents an instance of a name given to the EPUB Publication.
		Title []string `xml:"title"`
		// Specifies the language of the content of the given Rendition.
		Language []string `xml:"language"`
		// Represents the name of a person, organization, etc. responsible for the creation of the content of the Rendition
		Creator []struct {
			// The ID [XML] of the element, which MUST be unique within the document scope.
			ID string `xml:"id,attr,omitempty"`
			// Can be attached to the element to indicate the function the creator played in the creation of the content.
			Role   string `xml:"role,attr,omitempty"`
			FileAs string `xml:"file-as,attr,omitempty"`
			// Creator value
			Text string `xml:",chardata"`
		} `xml:"creator,omitempty"`
		// Description provides a description of the publication's content.
		Description []string `xml:"description,omitempty"`
		// Identifies the subject of the EPUB Publication.
		Subject []string `xml:"subject,omitempty"`
		// Identifies the publication's publisher.
		Publisher []string `xml:"publisher,omitempty"`
		// MUST only be used to define the publication date of the EPUB Publication.
		Date string `xml:"date"`
		// Provides a generic means of including package metadata.
		Meta []struct {
			// Takes a property data type value that defines the statement being made in the expression, and the text content of the element represents the assertion
			Property string `xml:"property,attr"`
			// Enhances the meaning of the expression or resource referenced
			Refines string `xml:"refines,attr,omitempty"`
			// Identifies metadata name - OPF2 extension
			Name string `xml:"name,attr"`
			// Contents metadata value of name  - OPF2 extension
			Content string `xml:"content,attr"`
			// Meta value
			Text string `xml:",chardata"`
		} `xml:"meta"`
	} `xml:"metadata"`
	// Provides an exhaustive list of the Publication Resources that constitute the given Rendition, each represented by an item element.
	Manifest struct {
		// Each item element in the manifest identifies a Publication Resource by the IRI provided in its href attribute.
		Item []struct {
			// The ID [XML] of the element, which MUST be unique within the document scope.
			ID string `xml:"id,attr"`
			// An absolute or relative IRI reference [RFC3987] to a resource.
			Href string `xml:"href,attr"`
			// Item element MUST conform to the applicable specification(s) as inferred from the MIME media type
			MediaType string `xml:"media-type,attr"`
			// A space-separated list of property values.
			Properties string `xml:"properties,attr,omitempty"`
		} `xml:"item"`
	} `xml:"manifest"`
}

func (opf *OPF) String() string {
	return "" + fmt.Sprint(
		"=========OPF===================\n",
		fmt.Sprintf("Lang:        %v\n", opf.Lang),
		fmt.Sprintf("Version:     %v\n", opf.Version),
		"---------Metadata--------------\n",
		fmt.Sprintf("Identifier:  %v\n", opf.Metadata.Identifier),
		fmt.Sprintf("Title:       %v\n", opf.Metadata.Title),
		fmt.Sprintf("Language:    %v\n", opf.Metadata.Language),
		fmt.Sprintf("Creator:     %v\n", opf.Metadata.Creator),
		fmt.Sprintf("Description: %v\n", opf.Metadata.Description),
		fmt.Sprintf("Subject:     %v\n", opf.Metadata.Subject),
		fmt.Sprintf("Publisher:   %v\n", opf.Metadata.Publisher),
		fmt.Sprintf("Date:        %v\n", opf.Metadata.Date),
		fmt.Sprintf("Meta:        %v\n", opf.Metadata.Meta),
		"---------Manifest--------------\n",
		fmt.Sprintf("Items:       %v\n", opf.Manifest.Item),
		"===============================\n",
	)
}

// Creates an opf package object from an OPF content.
func NewOPF(zr *zip.ReadCloser, path string) (*OPF, error) {
	r, err := zr.Open(path)
	if err != nil {
		return nil, err
	}
	opf := &OPF{}
	if err := decodeXML(r, &opf); err != nil {
		return nil, err
	}
	return opf, nil
}

func GetCoverImage(stock string, book *model.Book) (image.Image, error) {
	zr, _ := zip.OpenReader(path.Join(stock, book.File))
	defer zr.Close()
	var rc io.ReadCloser
	var err error
	for _, file := range zr.File {
		if strings.Contains(file.Name, book.Cover) {
			rc, err = file.Open()
			if err != nil {
				log.Fatal(err)
			}
			break
		}
	}
	defer rc.Close()
	img, _, err := image.Decode(bufio.NewReader(rc))

	if err != nil {
		return nil, err
	}
	return img, nil
}

// Utils ---------------------------------

func decodeXML(r io.Reader, v interface{}) error {
	decoder := xml.NewDecoder(r)
	decoder.Entity = xml.HTMLEntity
	decoder.CharsetReader = charset.NewReaderLabel
	return decoder.Decode(v)
}
