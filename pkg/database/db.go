package database

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/vinser/flibgolite/pkg/config"
	"github.com/vinser/flibgolite/pkg/hash"
	"github.com/vinser/flibgolite/pkg/model"
	"github.com/vinser/flibgolite/pkg/rlog"

	_ "embed"

	_ "modernc.org/sqlite"
)

const SQLITE_DB_BUSY_TIMEOUT = 10000

//go:embed sqlite_db_init.sql
var SQLITE_DB_INIT string

//go:embed sqlite_db_drop.sql
var SQLITE_DB_DROP string

// ==================================
type Handler struct {
	CFG    *config.Config
	Hashes *hash.BookHashes
	DB     *DB
	TX     *TX
	LOG    *rlog.Log
	WG     *sync.WaitGroup
	Queue  <-chan model.Book
	StopDB chan struct{}
}

type DB struct {
	*sqlx.DB
}

// ==================================
func NewDB(dsn string) *DB {
	err := os.MkdirAll(filepath.Dir(dsn), 0775)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
	// options := fmt.Sprintf("?_pragma=busy_timeout(%d)&_pragma=journal_mode(delete)", SQLITE_DB_BUSY_TIMEOUT)
	options := fmt.Sprintf("?_pragma=busy_timeout(%d)&_pragma=journal_mode(wal)", SQLITE_DB_BUSY_TIMEOUT)
	db := sqlx.MustOpen("sqlite", dsn+options)

	db.SetMaxOpenConns(30)
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	DB := &DB{
		DB: db,
	}

	return DB
}

func (db *DB) Close() {
	db.DB.Close()
}

func (db *DB) InitDB() {
	db.execFile(SQLITE_DB_INIT)
}

func (db *DB) DropDB() {
	if db.IsReady() {
		db.execFile(SQLITE_DB_DROP)
	}
}

func (db *DB) IsReady() bool {
	var err error
	rows, err := db.Query(`SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%' AND name NOT LIKE 'test%'`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	return rows.Next()
}

func (db *DB) execFile(sql string) {
	scanner := bufio.NewScanner(strings.NewReader(sql))
	scanner.Split(bufio.ScanLines)
	q := ""

	for scanner.Scan() {
		q += scanner.Text()
		if strings.Contains(q, ";") {
			_, err := db.Exec(q)
			q = ""
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

// ==================================
type TX struct {
	*sqlx.Tx
	Stmt map[string]*sqlx.Stmt
}

func (db *DB) txBegin() *TX {
	TX := &TX{
		Tx:   db.DB.MustBegin(),
		Stmt: map[string]*sqlx.Stmt{},
	}
	TX.PrepareStatements()
	return TX
}

func (tx *TX) txEnd() {
	defer func() {
		for _, stmt := range tx.Stmt {
			stmt.Close()
		}
		tx.Tx.Rollback()
	}()
	tx.Tx.Commit()
}

func (tx *TX) mustPrepare(query string) *sqlx.Stmt {
	stmt, err := tx.Tx.Preparex(query)
	if err != nil {
		panic(err)
	}
	return stmt
}
