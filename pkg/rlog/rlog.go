package rlog

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

type Log struct {
	*RotaryLog
	D *log.Logger // debug
	I *log.Logger // info
	W *log.Logger // warning
	E *log.Logger // error
}

func NewLog(logFile, level string) (Log *Log) {
	var w io.Writer
	if len(logFile) > 0 {
		err := os.MkdirAll(filepath.Dir(logFile), 0775)
		if err != nil && !os.IsExist(err) {
			log.Fatal(err)
		}
		w, err = NewRotaryLog(logFile, 0, 0, 0)
		// fw, err := NewRotaryLog(logFile, 60, 0, 0) // dedug 1 minute rotation
		if err != nil {
			log.Fatal(err)
		}
	} else {
		w = os.Stdout
	}
	d := io.Discard
	switch level {
	case "D":
		setLogWriters(Log, w, w, w, w)
	case "I":
		setLogWriters(Log, d, w, w, w)
	case "W":
		setLogWriters(Log, d, d, w, w)
	case "E":
		setLogWriters(Log, d, d, d, w)
	default:
	}
	return
}

func setLogWriters(l *Log, d, i, w, e io.Writer) {
	l.D = log.New(d, "DEBUG:   ", log.LstdFlags|log.Lshortfile)
	l.I = log.New(i, "INFO:    ", log.LstdFlags)
	l.W = log.New(w, "WARNING: ", log.LstdFlags)
	l.E = log.New(e, "ERROR:   ", log.LstdFlags|log.Lshortfile)
}
