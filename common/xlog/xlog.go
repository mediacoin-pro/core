package xlog

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

const (
	LevelNone    = 0
	LevelFatal   = 1
	LevelError   = 2
	LevelWarning = 3
	LevelInfo    = 4
	LevelDebug   = 5
	LevelTrace   = 6
)

var (
	Fatal   = newLogger("⛔ FATAL: ")
	Error   = newLogger("🛑 ERROR: ")
	Warning = newLogger("⚠️️ WARNG: ")
	Info    = newLogger("")
	Debug   = newLogger("ℹ️ Debug: ")
	Trace   = newLogger("-- Trace: ")
)

var (
	Level            = LevelWarning
	logOut io.Writer = os.Stderr
)

func init() {
	refreshLoggers()
}

func newLogger(prefix string) *log.Logger {
	return log.New(logOut, prefix, log.LstdFlags)
}

func refreshLoggers() {
	refreshLoger(Fatal, LevelFatal)
	refreshLoger(Error, LevelError)
	refreshLoger(Warning, LevelWarning)
	refreshLoger(Info, LevelInfo)
	refreshLoger(Debug, LevelDebug)
	refreshLoger(Trace, LevelTrace)
}

func refreshLoger(logger *log.Logger, loggerLv int) {
	if loggerLv <= Level {
		logger.SetOutput(logOut)
	} else {
		logger.SetOutput(ioutil.Discard)
	}
}

func SetLogLevel(lv int) {
	Level = lv
	refreshLoggers()
}

func SetOutput(w io.Writer) {
	logOut = w
	refreshLoggers()
}

func Panic(v ...interface{}) {
	Fatal.Panic(v...)
}

func Panicf(format string, v ...interface{}) {
	Fatal.Panicf(format, v...)
}

func Print(v ...interface{}) {
	Info.Print(v...)
}

func Printf(format string, v ...interface{}) {
	Info.Printf(format, v...)
}

func PrintfErr(format string, v ...interface{}) {
	if n := len(v); n > 0 {
		if err, _ := v[n-1].(error); err != nil {
			Error.Printf(format+" !!! ERROR: %v", v...)
		} else {
			Info.Printf(format, v[:n-1]...)
		}
	} else {
		Info.Print(format)
	}
}
