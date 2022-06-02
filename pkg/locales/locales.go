package locales

import (
	"errors"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gopkg.in/yaml.v3"
)

type Locales struct {
	DIR      string `yaml:"DIR"`
	DEFAULT  string `yaml:"DEFAULT"`
	ACCEPTED string `yaml:"ACCEPTED"`
	LANG     map[string]Language
}

type Language struct {
	Tag language.Tag
	Abc string
}

func (l *Locales) LoadLocales() {
	l.LANG = make(map[string]Language)
	var files []fs.DirEntry
	files, err := os.ReadDir(l.DIR)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			err := os.MkdirAll(l.DIR, 0775)
			if err != nil && !os.IsExist(err) {
				log.Fatal(err)
			}
			for lang, yml := range LOCALES_YML {
				err = os.WriteFile(filepath.Join(l.DIR, lang+".yml"), []byte(yml), 0775)
				if err != nil {
					log.Fatal(err)
				}
			}
		} else {
			log.Fatal(err)
		}
		files, _ = os.ReadDir(l.DIR)
	}

	for _, f := range files {
		if filepath.Ext(f.Name()) != ".yml" {
			continue
		}
		yamlFile, err := ioutil.ReadFile(filepath.Join(l.DIR, f.Name()))
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
				l.LANG[lang] = Language{Tag: lTag, Abc: splitABC(v)}
			default:
				message.SetString(lTag, k, v)
			}
		}
	}
}

func splitABC(abc string) string {
	s := "'"
	for _, r := range abc {
		s += string(r) + "', '"
	}
	return strings.TrimSuffix(s, ", '")
}
