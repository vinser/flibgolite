package rlog

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

type Log struct {
	File *RotaryLog
	I    *log.Logger // info
	E    *log.Logger // error
	D    *log.Logger // debug
}

func NewLog(logFile string, debug bool) *Log {
	dw := io.Discard
	err := os.MkdirAll(filepath.Dir(logFile), 0775)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
	// fw, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	fw, err := NewRotaryLog(logFile, 86400, 0, 0)
	// fw, err := NewRotaryLog(logFile, 60, 0, 0) // dedug 1 minute rotation
	if err != nil {
		log.Fatal(err)
	}
	if debug {
		dw = fw
	}
	return &Log{
		File: fw,
		I:    log.New(fw, "INFO:\t", log.LstdFlags),
		E:    log.New(fw, "ERROR:\t ", log.LstdFlags|log.Lshortfile),
		D:    log.New(dw, "DEBUG:\t", log.LstdFlags|log.Lshortfile),
	}
}
