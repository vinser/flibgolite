package rlog

import (
	"io/ioutil"
	"log"
)

type Log struct {
	File *RotaryLog
	I    *log.Logger // info
	E    *log.Logger // error
	D    *log.Logger // debug
}

func NewLog(logFile string, debug bool) *Log {
	dw := ioutil.Discard
	// fw, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	// fw, err := NewRotaryLog(logFile, 86400, 0, 0)
	fw, err := NewRotaryLog(logFile, 60, 0, 0)
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
