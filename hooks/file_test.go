// Copyright 2019 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hooks

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
)

var testLogDir = "/tmp/log"

func TestLogFile(t *testing.T) {
	logFileHook := LogFileHook{
		FileDir:  testLogDir,
		FileName: "spf.log",
	}
	logFileHook.Init()
	defer func() {
		os.RemoveAll(testLogDir)
	}()
	defer logFileHook.DeInit()
	l := logrus.New()
	l.AddHook(&logFileHook)
	l.Infof("hello, world")
	l.Debugf("hello, world?")
}

func TestLogFileRotate(t *testing.T) {
	logFn := "spf.rot.log"
	rotateNum := 10
	rotateSize := int64(1024)
	logFileHook := LogFileRotateHook{
		LogFileHook: LogFileHook{
			FileDir:  testLogDir,
			FileName: logFn,
		},
		RotateNum:  rotateNum,
		RotateSize: rotateSize,
	}
	{
		logFileHook.Init()
		defer func() {
			os.RemoveAll(testLogDir)
		}()
		defer logFileHook.DeInit()
		l := logrus.New()
		l.SetOutput(ioutil.Discard)
		l.AddHook(&logFileHook)
		for i := 0; i < 3828; i++ {
			l.Infof("%d", i)
		}
	}

	testSize := func(path string, ge bool, size int64) {
		fi, err := os.Lstat(path)
		if err != nil {
			if os.IsNotExist(err) {
				return
			}
			t.Errorf("%s: Lstat error: %s", path, err)
			return
		}
		got := fi.Size()
		if ge {
			if got >= size {
				return
			}
			t.Errorf("%s: got %d < size %d, expect >=", path, got, size)
		} else {
			if got < size {
				return
			}
			t.Errorf("%s: got %d >= size %d, expect <", path, got, size)
		}
		//return
	}
	fn := filepath.Join(testLogDir, logFn)
	testSize(fn, false, rotateSize)
	for i := 1; i < rotateNum; i++ {
		fn := filepath.Join(testLogDir, fmt.Sprintf("%s.%d", logFn, i))
		testSize(fn, true, rotateSize)
	}
}
