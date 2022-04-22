package config

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

// See config.yml for comments about this struct
type Config struct {
	Library struct {
		BOOK_STOCK       string `yaml:"BOOK_STOCK"`
		NEW_ACQUISITIONS string `yaml:"NEW_ACQUISITIONS"`
		TRASH            string `yaml:"TRASH"`
	}
	Language struct {
		LOCALES string `yaml:"LOCALES"`
		DEFAULT string `yaml:"DEFAULT"`
	}
	Database struct {
		DSN              string `yaml:"DSN"`
		INIT_SCRIPT      string `yaml:"INIT_SCRIPT"`
		DROP_SCRIPT      string `yaml:"DROP_SCRIPT"`
		POLL_DELAY       int    `yaml:"POLL_DELAY"`
		MAX_SCAN_THREADS int    `yaml:"MAX_SCAN_THREADS"`
		ACCEPTED_LANGS   string `yaml:"ACCEPTED_LANGS"`
	}
	Genres struct {
		TREE_FILE string `yaml:"TREE_FILE"`
	}
	Logs struct {
		OPDS  string `yaml:"OPDS"`
		SCAN  string `yaml:"SCAN"`
		DEBUG bool   `yaml:"DEBUG"`
	}
	OPDS struct {
		PORT      int `yaml:"PORT"`
		PAGE_SIZE int `yaml:"PAGE_SIZE"`
	}
}

func LoadConfig() *Config {
	var (
		b   []byte
		err error
	)
	expath, _ := os.Executable()
	dir, exname := filepath.Split(expath)
	ext := filepath.Ext(exname)

	configFile := filepath.Join(dir, "config", exname[:len(exname)-len(ext)]+".yml")

	b, err = os.ReadFile(configFile)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			err := os.MkdirAll(filepath.Dir(configFile), 0775)
			if err != nil && !os.IsExist(err) {
				log.Fatal(err)
			}
			err = os.WriteFile(configFile, []byte(CONFIG_YML), 0775)
			if err != nil {
				log.Fatal(err)
			}
			b, err = os.ReadFile(configFile)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	}
	c := &Config{}
	if err := yaml.Unmarshal([]byte(b), c); err != nil {
		log.Fatal(err)
	}
	return c
}

func (cfg *Config) LoadLocales() {
	files, err := os.ReadDir(cfg.Language.LOCALES)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			err := os.MkdirAll(cfg.Language.LOCALES, 0775)
			if err != nil && !os.IsExist(err) {
				log.Fatal(err)
			}
			err = os.WriteFile(filepath.Join(cfg.Language.LOCALES, "en.yml"), []byte(LOCALES_EN_YML), 0775)
			if err != nil {
				log.Fatal(err)
			}
			err = os.WriteFile(filepath.Join(cfg.Language.LOCALES, "ru.yml"), []byte(LOCALES_RU_YML), 0775)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	}

	for _, f := range files {
		if filepath.Ext(f.Name()) != ".yml" {
			continue
		}
		yamlFile, err := ioutil.ReadFile(filepath.Join(cfg.Language.LOCALES, f.Name()))
		if err != nil {
			log.Fatal(err)
		}
		data := map[string]string{}
		err = yaml.Unmarshal(yamlFile, &data)
		if err != nil {
			log.Fatal(err)
		}

		lang := language.Make(strings.TrimSuffix(f.Name(), ".yml"))
		for key, value := range data {
			message.SetString(lang, key, value)

		}
	}
}
