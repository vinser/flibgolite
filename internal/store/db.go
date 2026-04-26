package store

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	// "sync"

	"github.com/jmoiron/sqlx"

	_ "embed"

	_ "modernc.org/sqlite"
)

const SQLITE_DB_BUSY_TIMEOUT = 10000

//go:embed sqlite_db_init.sql
var SQLITE_DB_INIT string

//go:embed sqlite_db_drop.sql
var SQLITE_DB_DROP string

type DB struct {
	*sqlx.DB
}

// ==================================
func NewDB(dsn string) (*DB, error) {
	err := os.MkdirAll(filepath.Dir(dsn), 0775)
	if err != nil && !os.IsExist(err) {
		return nil, err
	}
	// options := fmt.Sprintf("?_pragma=busy_timeout(%d)&_pragma=journal_mode(delete)", SQLITE_DB_BUSY_TIMEOUT)
	options := fmt.Sprintf("?_pragma=busy_timeout(%d)&_pragma=journal_mode(wal)", SQLITE_DB_BUSY_TIMEOUT)
	db := sqlx.MustOpen("sqlite", dsn+options)

	db.SetMaxOpenConns(30)
	if err := db.Ping(); err != nil {
		return nil, err
	}

	DB := &DB{
		DB: db,
	}

	return DB, nil
}

func (db *DB) Close() {
	db.DB.Close()
}

func (db *DB) InitDB() {
	db.execFile(SQLITE_DB_INIT)
}

func (db *DB) DropDB() {
	ready, err := db.IsReady()
	if err != nil {
		// Handle error appropriately - for now we'll just skip dropping
		return
	}
	if ready {
		db.execFile(SQLITE_DB_DROP)
	}
}

func (db *DB) IsReady() (bool, error) {
	var err error
	rows, err := db.Query(`SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%' AND name NOT LIKE 'test%'`)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}

func (db *DB) execFile(sql string) error {
	scanner := bufio.NewScanner(strings.NewReader(sql))
	scanner.Split(bufio.ScanLines)
	q := ""

	for scanner.Scan() {
		q += scanner.Text()
		if strings.Contains(q, ";") {
			_, err := db.Exec(q)
			q = ""
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// ==================================
type TX struct {
	*sqlx.Tx
	Stmt map[string]*sqlx.Stmt
}

func (db *DB) TxBegin() *TX {
	TX := &TX{
		Tx:   db.DB.MustBegin(),
		Stmt: map[string]*sqlx.Stmt{},
	}
	TX.PrepareStatements()
	return TX
}

func (tx *TX) TxEnd() {
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
