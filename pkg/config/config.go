package config

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/vinser/flibgolite/pkg/locales"
	"github.com/vinser/flibgolite/pkg/rlog"
	"gopkg.in/yaml.v3"

	_ "embed"
)

// See config.yml for comments about this struct

type Library struct {
	STOCK_DIR string `yaml:"STOCK"`
	TRASH_DIR string `yaml:"TRASH"`
	NEW_DIR   string `yaml:"NEW"`
}
type Database struct {
	DSN               string `yaml:"DSN"`
	POLL_DELAY        int    `yaml:"POLL_DELAY"`
	MAX_SCAN_THREADS  int    `yaml:"MAX_SCAN_THREADS"`
	BOOK_QUEUE_SIZE   int    `yaml:"BOOK_QUEUE_SIZE"`
	FILE_QUEUE_SIZE   int    `yaml:"FILE_QUEUE_SIZE"`
	MAX_BOOKS_IN_TX   int    `yaml:"MAX_BOOKS_IN_TX"`
	DEDUPLICATE_LEVEL string `yaml:"DEDUPLICATE_LEVEL"`
}
type Genres struct {
	TREE_FILE string `yaml:"TREE_FILE"`
}
type Logs struct {
	OPDS  string `yaml:"OPDS"`
	SCAN  string `yaml:"SCAN"`
	LEVEL string `yaml:"LEVEL"`
}
type OPDS struct {
	PORT          int    `yaml:"PORT"`
	TITLE         string `yaml:"TITLE"`
	PAGE_SIZE     int    `yaml:"PAGE_SIZE"`
	LATEST_DAYS   int    `yaml:"LATEST_DAYS"`
	NO_CONVERSION bool   `yaml:"NO_CONVERSION"`
}
type Config struct {
	Library  Library  `yaml:"library"`
	Database Database `yaml:"database"`
	Genres   Genres   `yaml:"genres"`
	Logs     Logs     `yaml:"logs"`
	OPDS     OPDS     `yaml:"opds"`
	locales.Locales
}

func makeAbs(rootDir, path string) string {
	if path == "" {
		return ""
	}
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(rootDir, path)
}

//go:embed config.yml
var CONFIG_YML string

func LoadConfig(rootDir string) *Config {
	var (
		b   []byte
		err error
	)
	configFile := filepath.Join(rootDir, "config", "config.yml")

	b, err = os.ReadFile(configFile)
	if err != nil { // config file not found, create default
		if errors.Is(err, fs.ErrNotExist) {
			err := os.MkdirAll(filepath.Dir(configFile), 0775)
			if err != nil && !os.IsExist(err) {
				log.Fatal(err)
			}
			err = os.WriteFile(configFile, []byte(CONFIG_YML), 0664)
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
	c := &Config{
		Library: Library{
			STOCK_DIR: "books/stock",
			TRASH_DIR: "books/trash",
			NEW_DIR:   "",
		},
		Database: Database{
			DSN:               "dbdata/books.db",
			POLL_DELAY:        30,
			MAX_SCAN_THREADS:  10,
			BOOK_QUEUE_SIZE:   20000,
			FILE_QUEUE_SIZE:   20000,
			MAX_BOOKS_IN_TX:   20000,
			DEDUPLICATE_LEVEL: "F",
		},
		Genres: Genres{
			TREE_FILE: "config/genres.xml",
		},
		Logs: Logs{
			OPDS:  "logs/opds.log",
			SCAN:  "logs/scan.log",
			LEVEL: "W",
		},
		OPDS: OPDS{
			PORT:          8085,
			TITLE:         "FLib Go Go Go!!!",
			PAGE_SIZE:     30,
			LATEST_DAYS:   14,
			NO_CONVERSION: false,
		},
		Locales: locales.Locales{
			DIR:      "config/locales",
			DEFAULT:  "en",
			ACCEPTED: "en, ru, uk",
		},
	}
	if err := yaml.Unmarshal([]byte(b), c); err != nil {
		log.Fatal(err)
	}

	c.Library.STOCK_DIR = makeAbs(rootDir, c.Library.STOCK_DIR)
	if len(c.Library.TRASH_DIR) > 0 {
		c.Library.TRASH_DIR = makeAbs(rootDir, c.Library.TRASH_DIR)
	}
	if len(c.Library.NEW_DIR) > 0 {
		c.Library.NEW_DIR = makeAbs(rootDir, c.Library.NEW_DIR)
	}
	c.Locales.DIR = makeAbs(rootDir, c.Locales.DIR)
	c.Genres.TREE_FILE = makeAbs(rootDir, c.Genres.TREE_FILE)
	c.Database.DSN = makeAbs(rootDir, c.Database.DSN)
	c.Logs.OPDS = makeAbs(rootDir, c.Logs.OPDS)
	c.Logs.SCAN = makeAbs(rootDir, c.Logs.SCAN)

	return c
}

func (c *Config) InitLogs(needOpds bool) (stockLog, opdsLog *rlog.Log) {
	stockLog = nil
	opdsLog = nil
	stockLog = rlog.NewLog(c.Logs.SCAN, c.Logs.LEVEL)
	if needOpds {
		if c.Logs.SCAN == c.Logs.OPDS {
			opdsLog = stockLog
		} else {
			opdsLog = rlog.NewLog(c.Logs.OPDS, c.Logs.LEVEL)
		}
	}
	return
}
