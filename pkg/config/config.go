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
)

// See config.yml for comments about this struct
type Config struct {
	Library struct {
		STOCK_DIR string `yaml:"STOCK"`
		NEW_DIR   string `yaml:"NEW"`
		TRASH_DIR string `yaml:"TRASH"`
	}
	Database struct {
		DSN              string `yaml:"DSN"`
		INIT_SCRIPT      string `yaml:"INIT_SCRIPT"`
		DROP_SCRIPT      string `yaml:"DROP_SCRIPT"`
		POLL_DELAY       int    `yaml:"POLL_DELAY"`
		MAX_SCAN_THREADS int    `yaml:"MAX_SCAN_THREADS"`
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
	locales.Locales
}

func makeAbs(path string) string {
	if path == "" {
		return ""
	}
	if filepath.IsAbs(path) {
		return path
	}
	execpath, err := os.Executable()
	if err != nil {
		return "."
	}
	return filepath.Join(filepath.Dir(execpath), path)
}

func LoadConfig() *Config {
	var (
		b   []byte
		err error
	)
	expath, _ := os.Executable()
	dir := filepath.Dir(expath)
	// dir, exname := filepath.Split(expath)
	// ext := filepath.Ext(exname)

	// configFile := filepath.Join(dir, "config", exname[:len(exname)-len(ext)]+".yml")
	configFile := filepath.Join(dir, "config", "config.yml")

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

	c.Library.STOCK_DIR = makeAbs(c.Library.STOCK_DIR)
	c.Library.NEW_DIR = makeAbs(c.Library.NEW_DIR)
	c.Library.TRASH_DIR = makeAbs(c.Library.TRASH_DIR)
	c.Locales.DIR = makeAbs(c.Locales.DIR)
	c.Genres.TREE_FILE = makeAbs(c.Genres.TREE_FILE)
	c.Database.DSN = makeAbs(c.Database.DSN)
	c.Logs.OPDS = makeAbs(c.Logs.OPDS)
	c.Logs.SCAN = makeAbs(c.Logs.SCAN)

	return c
}

func (c *Config) InitLogs(needOpds bool) (stockLog, opdsLog *rlog.Log) {
	stockLog = nil
	opdsLog = nil
	stockLog = rlog.NewLog(c.Logs.SCAN, c.Logs.DEBUG)
	if needOpds {
		if c.Logs.SCAN == c.Logs.OPDS {
			opdsLog = stockLog
		} else {
			opdsLog = rlog.NewLog(c.Logs.OPDS, c.Logs.DEBUG)
		}
	}
	return
}
