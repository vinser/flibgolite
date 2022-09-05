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
	S *log.Logger // service
}

func NewLog(logFile, level string) *Log {
	var w io.Writer
	l := &Log{}
	if len(logFile) > 0 {
		err := os.MkdirAll(filepath.Dir(logFile), 0775)
		if err != nil && !os.IsExist(err) {
			log.Fatal(err)
		}
		rl, err := NewRotaryLog(logFile, 0, 0, 0)
		if err != nil {
			log.Fatal(err)
		}
		l.RotaryLog = rl
		w = rl

	} else {
		w = os.Stdout
	}
	d := io.Discard
	switch level {
	case "D":
		setLogWriters(l, w, w, w, w)
	case "I":
		setLogWriters(l, d, w, w, w)
	case "W":
		setLogWriters(l, d, d, w, w)
	case "E", "S": // Always
		setLogWriters(l, d, d, d, w)
	default:
		setLogWriters(l, d, d, w, w)
	}
	return l
}

func setLogWriters(l *Log, d, i, w, e io.Writer) {
	l.D = log.New(d, "DEBUG:   ", log.LstdFlags|log.Lshortfile)
	l.I = log.New(i, "INFO:    ", log.LstdFlags)
	l.W = log.New(w, "WARNING: ", log.LstdFlags)
	l.E = log.New(e, "ERROR:   ", log.LstdFlags|log.Lshortfile)
	l.S = log.New(e, "SERVICE: ", log.LstdFlags)
}
