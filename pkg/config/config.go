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
type Config struct {
	Library struct {
		STOCK_DIR string `yaml:"STOCK"`
		TRASH_DIR string `yaml:"TRASH"`
		NEW_DIR   string `yaml:"NEW"`
	}
	Database struct {
		DSN               string `yaml:"DSN"`
		POLL_DELAY        int    `yaml:"POLL_DELAY"`
		MAX_SCAN_THREADS  int    `yaml:"MAX_SCAN_THREADS"`
		BOOK_QUEUE_SIZE   int    `yaml:"BOOK_QUEUE_SIZE"`
		FILE_QUEUE_SIZE   int    `yaml:"FILE_QUEUE_SIZE"`
		MAX_BOOKS_IN_TX   int    `yaml:"MAX_BOOKS_IN_TX"`
		DEDUPLICATE_LEVEL string `yaml:"DEDUPLICATE_LEVEL"`
	}
	Genres struct {
		TREE_FILE string `yaml:"TREE_FILE"`
	}
	Logs struct {
		OPDS  string `yaml:"OPDS"`
		SCAN  string `yaml:"SCAN"`
		LEVEL string `yaml:"LEVEL"`
	}
	OPDS struct {
		PORT          int    `yaml:"PORT"`
		TITLE         string `yaml:"TITLE"`
		PAGE_SIZE     int    `yaml:"PAGE_SIZE"`
		LATEST_DAYS   int    `yaml:"LATEST_DAYS"`
		NO_CONVERSION bool   `yaml:"NO_CONVERSION"`
	}
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
	c := &Config{}
	if err := yaml.Unmarshal([]byte(b), c); err != nil {
		log.Fatal(err)
	}

	if c.Library.STOCK_DIR == "" {
		c.Library.STOCK_DIR = "books/stock"
	}
	if c.Library.TRASH_DIR == "" {
		c.Library.TRASH_DIR = "books/trash"
	}
	if c.Database.DSN == "" {
		c.Database.DSN = "dbdata/books.db"
	}

	if c.Genres.TREE_FILE == "" {
		c.Genres.TREE_FILE = "config/genres.xml"
	}

	if c.Database.POLL_DELAY == 0 {
		c.Database.POLL_DELAY = 30
	}
	if c.Database.MAX_SCAN_THREADS == 0 {
		c.Database.MAX_SCAN_THREADS = 10
	}
	if c.Database.BOOK_QUEUE_SIZE == 0 {
		c.Database.BOOK_QUEUE_SIZE = 20000
	}
	if c.Database.FILE_QUEUE_SIZE == 0 {
		c.Database.FILE_QUEUE_SIZE = 20000
	}
	if c.Database.MAX_BOOKS_IN_TX == 0 {
		c.Database.MAX_BOOKS_IN_TX = 20000
	}
	if c.Database.DEDUPLICATE_LEVEL == "" {
		c.Database.DEDUPLICATE_LEVEL = "F"
	}

	if c.Logs.OPDS == "" {
		c.Logs.OPDS = "logs/opds.log"
	}
	if c.Logs.SCAN == "" {
		c.Logs.SCAN = "logs/scan.log"
	}
	if c.Logs.LEVEL == "" {
		c.Logs.LEVEL = "W"
	}

	if c.OPDS.PORT == 0 {
		c.OPDS.PORT = 8085
	}
	if c.OPDS.TITLE == "" {
		c.OPDS.TITLE = "FLib Go Go Go!!!"
	}
	if c.OPDS.PAGE_SIZE == 0 {
		c.OPDS.PAGE_SIZE = 30
	}
	if c.OPDS.LATEST_DAYS == 0 {
		c.OPDS.LATEST_DAYS = 14
	}

	if c.Locales.DIR == "" {
		c.Locales.DIR = "config/locales"
	}
	if c.Locales.DEFAULT == "" {
		c.Locales.DEFAULT = "en"
	}
	if c.Locales.ACCEPTED == "" {
		c.Locales.ACCEPTED = "en, ru, uk"
	}

	c.Library.STOCK_DIR = makeAbs(rootDir, c.Library.STOCK_DIR)
	c.Library.TRASH_DIR = makeAbs(rootDir, c.Library.TRASH_DIR)
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
