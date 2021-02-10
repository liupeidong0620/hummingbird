package log

import (
	"fmt"
	"io"
	"log"
	"os"
)

const (
	LogLevelTrace = iota
	LogLevelDebug
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

var logger = log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile|log.LUTC|log.Lmicroseconds)
var level = LogLevelInfo

func SetLevel(lev int) {
	level = lev
}

func SetOutput(w io.Writer) {
	logger.SetOutput(w)
}

func SetPrefix(prefix string) {
	logger.SetPrefix(prefix)
}

func Trace(v ...interface{}) {
	if level <= LogLevelTrace {
		msg := fmt.Sprintln(v...)
		logger.Output(2, fmt.Sprintf("[TRACE] %s", msg))
	}
}

func Debug(v ...interface{}) {
	if level <= LogLevelDebug {
		msg := fmt.Sprintln(v...)
		logger.Output(2, fmt.Sprintf("[DEBUG] %s", msg))
	}
}

func Info(v ...interface{}) {
	if level <= LogLevelInfo {
		msg := fmt.Sprintln(v...)
		logger.Output(2, fmt.Sprintf("[INFO] %s", msg))
	}
}

func Warn(v ...interface{}) {
	if level <= LogLevelWarn {
		msg := fmt.Sprintln(v...)
		logger.Output(2, fmt.Sprintf("[WARN] %s", msg))
	}
}

func Error(v ...interface{}) {
	if level <= LogLevelError {
		msg := fmt.Sprintln(v...)
		logger.Output(2, fmt.Sprintf("[ERROR] %s", msg))
	}
}
