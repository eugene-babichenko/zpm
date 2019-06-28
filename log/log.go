package log

import (
	"fmt"
	"os"
)

type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	ErrorLevel
	FatalLevel
)

var levelUsed Level

func SetLevel(level Level) {
	switch level {
	case DebugLevel, InfoLevel, ErrorLevel, FatalLevel:
		levelUsed = level
	}
}

func writeMsg(level Level, format string, args ...interface{}) {
	if levelUsed > level {
		return
	}

	var str string
	if len(args) > 0 {
		str = fmt.Sprintf(format, args...)
	} else {
		str = format
	}

	var levelName string
	switch level {
	case DebugLevel:
		levelName = "debug"
	case InfoLevel:
		levelName = "info"
	case ErrorLevel:
		levelName = "error"
	case FatalLevel:
		levelName = "fatal"
	}

	_, _ = fmt.Fprintf(os.Stderr, "%s: %s\n", levelName, str)
}

func Debug(format string, args ...interface{}) {
	writeMsg(DebugLevel, format, args...)
}

func Info(format string, args ...interface{}) {
	writeMsg(InfoLevel, format, args...)
}

func Error(format string, args ...interface{}) {
	writeMsg(ErrorLevel, format, args...)
}

func Fatal(format string, args ...interface{}) {
	writeMsg(FatalLevel, format, args...)
	os.Exit(1)
}
