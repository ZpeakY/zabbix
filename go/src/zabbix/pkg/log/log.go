/*
** Zabbix
** Copyright (C) 2001-2019 Zabbix SIA
**
** This program is free software; you can redistribute it and/or modify
** it under the terms of the GNU General Public License as published by
** the Free Software Foundation; either version 2 of the License, or
** (at your option) any later version.
**
** This program is distributed in the hope that it will be useful,
** but WITHOUT ANY WARRANTY; without even the implied warranty of
** MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
** GNU General Public License for more details.
**
** You should have received a copy of the GNU General Public License
** along with this program; if not, write to the Free Software
** Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
**/

package log

import (
	"bytes"
	"errors"
	"log"
	"os"
	"runtime/debug"
)

const Empty = 0
const Crit = 1
const Err = 2
const Warning = 3
const Debug = 4
const Trace = 5

const Info = 127

const Undefined = 0
const System = 1
const File = 2
const Console = 3

var logLevel int
var logger *log.Logger

func CheckLogLevel(level int) bool {
	if level != Info && (level > logLevel || Empty == level) {
		return false
	}
	return true
}

func IncreaseLogLevel() {
	if logLevel != Trace {
		logLevel++
	}
}

func DecreaseLogLevel() {
	if logLevel != Empty {
		logLevel--
	}
}

func Open(logType int, level int, filename string) error {

	if logType == Console {
		logger = log.New(os.Stdout, "", log.Lmicroseconds|log.Ldate)
	} else if logType == File {
		f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		logger = log.New(f, "", log.Lmicroseconds|log.Ldate)
	} else {
		return errors.New("invalid argument")
	}

	logLevel = level
	return nil
}

func Critf(format string, args ...interface{}) {
	if CheckLogLevel(Crit) {
		logger.Printf(format, args...)
	}
}

func Infof(format string, args ...interface{}) {
	logger.Printf(format, args...)
}

func Warningf(format string, args ...interface{}) {
	if CheckLogLevel(Warning) {
		logger.Printf(format, args...)
	}
}

func Tracef(format string, args ...interface{}) {
	if CheckLogLevel(Trace) {
		logger.Printf(format, args...)
	}
}

func Debugf(format string, args ...interface{}) {
	if CheckLogLevel(Debug) {
		logger.Printf(format, args...)
	}
}

func Errf(format string, args ...interface{}) {
	if CheckLogLevel(Err) {
		logger.Printf(format, args...)
	}
}

func PanicHook() {
	if r := recover(); r != nil {
		data := debug.Stack()
		Critf("Critical failure: %v", r)
		var tail int
		for offset, end, num := 0, 0, 1; end != -1; offset, num = offset+end+1, num+1 {
			end = bytes.IndexByte(data[offset:], '\n')
			if end != -1 {
				tail = offset + end
			} else {
				tail = len(data)
			}
			Critf("%s", string(data[offset:tail]))
		}
		panic(r)
	}
}
