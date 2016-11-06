package log

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

const (
	//can be refer outside
	DEBUG = 0
	INFO  = 1
	WARN  = 2
	ERROR = 3

	prefixDebug = "[DEBUG] "
	prefixInfo  = "[INFO ] "
	prefixWarn  = "[WARN ] "
	prefixError = "[ERROR] "
)

type Logger struct {
	level    int
	goLogger *log.Logger
	baseFile *os.File
}

func New(level int, filename string) (*Logger, error) {
	if level > ERROR || level < DEBUG {
		return nil, errors.New("not a valid log level")
	}

	flag := log.LstdFlags /* | log.Lshortfile */

	var tLog *log.Logger
	if filename != "" {
		dir := filename[:strings.LastIndex(filename, "/")]
		if !existPath(dir) {
			os.MkdirAll(dir, os.ModePerm)
		}
		var file *os.File
		if !existPath(filename) {
			f, err := os.Create(filename)
			if err != nil {
				return nil, err
			}
			file = f
		} else {
			f, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0666)
			if err != nil {
				return nil, err
			}
			file = f
		}

		tLog = log.New(file, "", flag)
	} else {
		tLog = log.New(os.Stdout, "", flag)
	}

	logger := &Logger{}
	logger.level = level
	logger.goLogger = tLog
	return logger, nil
}

func existPath(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
func dayString() string {
	now := time.Now()
	f := fmt.Sprintf("%d%02d%02d",
		now.Year(),
		now.Month(),
		now.Day(),
	)
	return f
}

func (logger *Logger) Close() {
	if logger.baseFile != nil {
		logger.baseFile.Close()
	}

	logger.goLogger = nil
	logger.baseFile = nil
}

func (logger *Logger) print0(level int, levelStr string, format string, a ...interface{}) {
	if level < logger.level {
		return
	}
	if logger.goLogger == nil {
		panic("logger closed")
	}

	format = levelStr + format
	logger.goLogger.Printf(format, a...)
}

func (logger *Logger) Debug(format string, a ...interface{}) {
	logger.print0(DEBUG, prefixDebug, format, a...)
}

func (logger *Logger) Info(format string, a ...interface{}) {
	logger.print0(INFO, prefixInfo, format, a...)
}

func (logger *Logger) Warn(format string, a ...interface{}) {
	logger.print0(WARN, prefixWarn, format, a...)
}

func (logger *Logger) Error(format string, a ...interface{}) {
	logger.print0(ERROR, prefixError, format, a...)
}
func (logger *Logger) Err(err error, format string, a ...interface{}) {
	logger.Error(format+"\n %v\n%v", append(a, err.Error(), string(debug.Stack()))...)
}
