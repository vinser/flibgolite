/*
Copyright (c) 2016, immortal
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

* Neither the name of logrotate nor the names of its
  contributors may be used to endorse or promote products derived from
  this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/
package rlog

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var _ io.WriteCloser = (*RotaryLog)(nil)

// RotaryLog struct
type RotaryLog struct {
	sync.Mutex
	Age   time.Duration
	Num   int
	Size  int
	file  *os.File
	sTime time.Time
	size  int64
}

// New return instance of RotaryLog
// defaults
// age  86400 rotate every 24h0m0s
// num  7     files
// size 0     no limit size
func NewRotaryLog(logfile string, age, num, size int) (*RotaryLog, error) {
	f, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	Age := time.Duration(0)
	if age > 0 {
		Age = time.Duration(age) * time.Second
	}
	num--
	if num < 0 {
		num = 7
	}
	Size := 0
	if size > 0 {
		Size = size * 1048576
	}

	rl := &RotaryLog{
		Mutex: sync.Mutex{},
		Age:   Age,
		Num:   num,
		Size:  Size,
		file:  f,
		sTime: time.Now(),
		size:  0,
	}
	i, err := rl.file.Stat()
	if err != nil {
		return rl, nil
	}
	rl.sTime = i.ModTime()
	rl.size = i.Size()
	// rotate if needed
	if rl.Age > 0 && time.Since(rl.sTime) >= rl.Age || rl.Size > 0 && rl.size > int64(rl.Size) {
		if err := rl.rotate(); err != nil {
			return nil, err
		}
	}
	return rl, nil
}

// Write implements io.Writer
func (rl *RotaryLog) Write(p []byte) (n int, err error) {
	rl.Lock()
	defer rl.Unlock()

	writeLen := int64(len(p))

	// rotate based on Age and Size
	if rl.Age > 0 && time.Since(rl.sTime) >= rl.Age || rl.Size > 0 && rl.size+writeLen > int64(rl.Size) {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = rl.file.Write(p)
	rl.size += int64(n)
	return n, err
}

// Close implements io.Closer, and closes the current logfile
func (rl *RotaryLog) Close() error {
	rl.Lock()
	defer rl.Unlock()
	return rl.close()
}

// close closes the file if it is open
func (rl *RotaryLog) close() error {
	if rl.file == nil {
		return nil
	}
	err := rl.file.Close()
	rl.file = nil
	return err
}

// Rotate helper function for rotate
func (rl *RotaryLog) Rotate() error {
	rl.Lock()
	defer rl.Unlock()
	return rl.rotate()
}

// rotate close existing log file and create a new one
func (rl *RotaryLog) rotate() error {
	digits := fmt.Sprint(len(strconv.Itoa(rl.Num)))
	format := "%s-%0" + digits + "d%s"
	path := rl.file.Name()
	ext := filepath.Ext(path)
	rl.close()
	// rotate logs
	for i := rl.Num; i >= 0; i-- {
		logfile := fmt.Sprintf(format, strings.TrimSuffix(path, ext), i, ext)
		if _, err := os.Stat(logfile); err == nil {
			// delete old file
			if i == rl.Num {
				os.Remove(logfile)
			} else if err := os.Rename(logfile, fmt.Sprintf(format, strings.TrimSuffix(path, ext), i+1, ext)); err != nil {
				return err
			}
		}
	}
	// create logfile 0
	if err := os.Rename(path, fmt.Sprintf(format, strings.TrimSuffix(path, ext), 0, ext)); err != nil {
		return err
	}
	// create new log file
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	rl.file = f
	rl.sTime = time.Now()
	rl.size = 0
	return nil
}
