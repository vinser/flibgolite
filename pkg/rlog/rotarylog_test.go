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
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewRotaryLog(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "TestNew")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	type Expected struct {
		Age       time.Duration
		Num, Size int
	}

	var testArgs = []struct {
		args     []int
		expected Expected
	}{
		{[]int{0, 0, 0}, Expected{time.Duration(0), 7, 0}},
		{[]int{0, 0, 1}, Expected{time.Duration(0), 7, 1048576}},
		{[]int{0, 1, 0}, Expected{time.Duration(0), 0, 0}},
		{[]int{0, 1, 1}, Expected{time.Duration(0), 0, 1048576}},
		{[]int{1, 0, 0}, Expected{time.Duration(1) * time.Second, 7, 0}},
		{[]int{1, 0, 1}, Expected{time.Duration(1) * time.Second, 7, 1048576}},
		{[]int{1, 1, 1}, Expected{time.Duration(1) * time.Second, 0, 1048576}},
		{[]int{0, 3, 1}, Expected{time.Duration(0), 2, 1048576}},
		{[]int{86400, 0, 1}, Expected{time.Duration(86400) * time.Second, 7, 1048576}},
		{[]int{43200, 0, 1}, Expected{time.Duration(43200) * time.Second, 7, 1048576}},
	}

	for _, a := range testArgs {
		l, err := NewRotaryLog(tmpfile.Name(), a.args[0], a.args[1], a.args[2])
		if err != nil {
			t.Error(err)
		}
		if l.Age != a.expected.Age {
			t.Errorf("Expecting age %v, got: %v", a.expected.Age, l.Age)
		}
		if l.Num != a.expected.Num {
			t.Errorf("Expecting num %v, got: %v", a.expected.Num, l.Num)
		}
		if l.Size != a.expected.Size {
			t.Errorf("Expecting size %v, got: %v", a.expected.Size, l.Size)
		}
	}
}

func TestRotate(t *testing.T) {
	var testRotate = []struct {
		args     []int
		expected int
	}{
		{[]int{0, 0, 0}, 1},
		{[]int{0, 0, 1}, 1},
		{[]int{1, 1, 0}, 2},
		{[]int{1, 0, 0}, 4},
		{[]int{1, 3, 0}, 4},
	}

	for _, a := range testRotate {
		dir, err := ioutil.TempDir("", "TestRotate")
		if err != nil {
			t.Error(err)
		}
		tmplog := filepath.Join(dir, "test.log")
		l, err := NewRotaryLog(tmplog, a.args[0], a.args[1], a.args[2])
		if err != nil {
			t.Error(err)
		}
		log.SetOutput(l)
		for i := 0; i <= 5; i++ {
			time.Sleep(500 * time.Millisecond)
			log.Println(i)
		}
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			t.Fatal(err)
		}
		if len(files) != a.expected {
			os.RemoveAll(dir)
			t.Fatalf("Expecting %v got %v", a.expected, len(files))
		}
		os.RemoveAll(dir)
	}
}

func TestRotateRotate(t *testing.T) {
	dir, err := ioutil.TempDir("", "TestRotateRotate")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	tmplog := filepath.Join(dir, "test.log")
	d1 := []byte("not\nempty\n")
	err = ioutil.WriteFile(tmplog, d1, 0644)
	if err != nil {
		t.Error(err)
	}
	l, err := NewRotaryLog(tmplog, 0, 0, 0)
	if err != nil {
		t.Error(err)
	}
	l.Rotate()
	log.SetOutput(l)
	for i := 0; i <= 100; i++ {
		log.Println(i)
	}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 2 {
		t.Errorf("Expecting 2 files got: %v", len(files))
	}
	l.Rotate()
	files, err = ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 3 {
		t.Errorf("Expecting 3 files got: %v", len(files))
	}
}

func TestNewRotateAge(t *testing.T) {
	dir, err := ioutil.TempDir("", "TestRotateAge")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	tmplog := filepath.Join(dir, "test.log")
	// fmt.Println("tmplog: ", tmplog)
	d1 := []byte("not\nempty\n")
	err = ioutil.WriteFile(tmplog, d1, 0644)
	if err != nil {
		t.Error(err)
	}
	// myTime, _ := time.Parse(time.RFC822, "01 Jan 01 00:00 UTC")
	myTime := time.Now().Add(-time.Second * 86400)
	err = os.Chtimes(tmplog, myTime, myTime)
	if err != nil {
		fmt.Println(err)
	}
	_, err = NewRotaryLog(tmplog, 86400, 0, 0)
	if err != nil {
		t.Error(err)
	}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 2 {
		t.Errorf("Expecting 2 files got: %v", len(files))
	}
}

func TestNewRotateSize(t *testing.T) {
	dir, err := ioutil.TempDir("", "TestRotateSize")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	tmplog := filepath.Join(dir, "test.log")
	d1 := []byte("not\nempty\n")
	err = ioutil.WriteFile(tmplog, d1, 0644)
	if err != nil {
		t.Error(err)
	}
	err = os.Truncate(tmplog, 1048577)
	if err != nil {
		fmt.Println(err)
	}
	_, err = NewRotaryLog(tmplog, 0, 0, 1)
	if err != nil {
		t.Error(err)
	}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 2 {
		t.Errorf("Expecting 2 files got: %v", len(files))
	}
}
