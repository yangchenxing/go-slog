package slog

import (
	"fmt"
	"runtime/debug"
)

var (
	globalFields  = make(map[string]interface{})
	GlobalSession = Session(make(map[string]interface{}))
)

func WithField(key string, value interface{}) {
	globalFields[key] = value
}

func WithFields(fields map[string]interface{}) {
	for key, value := range fields {
		globalFields[key] = value
	}
}

func Debug(args ...interface{}) {
	GlobalSession.EventSkip(2).Log("debug", args...)
}

func Debugf(format string, args ...interface{}) {
	GlobalSession.EventSkip(2).Logf("debug", format, args...)
}

func Debugln(args ...interface{}) {
	GlobalSession.EventSkip(2).Logln("debug", args...)
}

func Info(args ...interface{}) {
	GlobalSession.EventSkip(2).Log("info", args...)
}

func Infof(format string, args ...interface{}) {
	GlobalSession.EventSkip(2).Logf("info", format, args...)
}

func Infoln(args ...interface{}) {
	GlobalSession.EventSkip(2).Logln("info", args...)
}

func Warn(args ...interface{}) {
	GlobalSession.EventSkip(2).Log("warn", args...)
}

func Warnf(format string, args ...interface{}) {
	GlobalSession.EventSkip(2).Logf("warn", format, args...)
}

func Warnln(args ...interface{}) {
	GlobalSession.EventSkip(2).Logln("warn", args...)
}

func Error(args ...interface{}) {
	GlobalSession.EventSkip(2).Log("error", args...)
}

func Errorf(format string, args ...interface{}) {
	GlobalSession.EventSkip(2).Logf("error", format, args...)
}

func Errorln(args ...interface{}) {
	GlobalSession.EventSkip(2).Logln("error", args...)
}

func Fatal(args ...interface{}) {
	GlobalSession.EventSkip(2).Log("fatal", args...)
}

func Fatalf(format string, args ...interface{}) {
	GlobalSession.EventSkip(2).Logf("fatal", format, args...)
}

func Fatalln(args ...interface{}) {
	GlobalSession.EventSkip(2).Logln("fatal", args...)
}

func Panic(args ...interface{}) {
	event := GlobalSession.EventSkip(2)
	event.Level = "panic"
	event.Message = fmt.Sprint(args...)
	event.WithField(stackKey, string(debug.Stack()))
	panic(event)
}

func Panicf(format string, args ...interface{}) {
	event := GlobalSession.EventSkip(2)
	event.Level = "panic"
	event.Message = fmt.Sprintf(format, args...)
	event.WithField(stackKey, string(debug.Stack()))
	panic(event)
}

func Panicln(args ...interface{}) {
	event := GlobalSession.EventSkip(2)
	event.Level = "panic"
	event.Message = fmt.Sprintln(args...)
	event.WithField(stackKey, string(debug.Stack()))
	panic(event)
}
