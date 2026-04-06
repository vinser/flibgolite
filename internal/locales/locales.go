package locales

import (
	"embed"
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gopkg.in/yaml.v3"
)

type Locales struct {
	DIR       string `yaml:"DIR"`
	DEFAULT   string `yaml:"DEFAULT"`
	ACCEPTED  string `yaml:"ACCEPTED"`
	Languages map[string]Language
	Matcher   language.Matcher
}

type Language struct {
	Tag language.Tag
	Abc string
}

//go:embed *.yml
var LOCALES_YML embed.FS

func (l *Locales) LoadLocales() {
	l.Languages = make(map[string]Language)

	if err := os.MkdirAll(l.DIR, 0775); err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	defaultFile := l.DEFAULT + ".yml"
	defaultFilePath := filepath.Join(l.DIR, defaultFile)

	if _, err := os.Stat(defaultFilePath); errors.Is(err, fs.ErrNotExist) {
		if _, err := LOCALES_YML.ReadFile(defaultFile); err != nil {
			log.Fatal("default locale file not found: " + defaultFilePath)
		}

		ymls, err := LOCALES_YML.ReadDir(".")
		if err != nil {
			log.Fatal(err)
		}
		for _, yml := range ymls {
			src, err := LOCALES_YML.ReadFile(yml.Name())
			if err != nil {
				log.Fatal(err)
			}
			err = os.WriteFile(filepath.Join(l.DIR, yml.Name()), src, 0664)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	files, err := os.ReadDir(l.DIR)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if filepath.Ext(f.Name()) != ".yml" {
			continue
		}
		yamlFile, err := os.ReadFile(filepath.Join(l.DIR, f.Name()))
		if err != nil {
			log.Fatal(err)
		}
		data := map[string]string{}
		err = yaml.Unmarshal(yamlFile, &data)
		if err != nil {
			log.Fatal(err)
		}

		lang := strings.TrimSuffix(f.Name(), ".yml")
		lTag := language.Make(lang)
		for k, v := range data {
			switch k {
			case "ABC":
				l.Languages[lang] = Language{Tag: lTag, Abc: splitABC(v)}
			default:
				message.SetString(lTag, k, v)
			}
		}
	}

	l.newMatcher()
}

func (l *Locales) newMatcher() {
	tags := []language.Tag{}
	defTag := language.Make(l.DEFAULT)
	tags = append(tags, defTag)
	for _, lang := range l.Languages {
		if lang.Tag != defTag {
			tags = append(tags, lang.Tag)
		}
	}
	l.Matcher = language.NewMatcher(tags)
}

func splitABC(abc string) string {
	s := "'"
	for _, r := range abc {
		s += string(r) + "', '"
	}
	return strings.TrimSuffix(s, ", '")
}
