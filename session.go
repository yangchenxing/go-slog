package slog

import (
	"fmt"
	"runtime/debug"
)

// Session stores special fields for a session in app
type Session map[string]interface{}

// NewSession create a new session instance
func NewSession() Session {
	return Session(make(map[string]interface{}))
}

// WithField add a key value pair to session
func (session Session) WithField(key string, value interface{}) Session {
	session[key] = value
	return session
}

// WithFields add multiple key value pairs to session
func (session Session) WithFields(fields Fields) Session {
	for key, value := range fields {
		session[key] = value
	}
	return session
}

// Event create a new event of session
func (session Session) Event() *Event {
	return session.EventSkip(2)
}

// EventSkip create a new event with special skip
func (session Session) EventSkip(skip int) *Event {
	return newEvent(skip+1, session)
}

// Debug log `debug` event with fmt.Sprint
func (session Session) Debug(args ...interface{}) {
	session.EventSkip(2).Log(debugLevel, args...)
}

// Debugf log `debug` event with fmt.Sprintf
func (session Session) Debugf(format string, args ...interface{}) {
	session.EventSkip(2).Logf(debugLevel, format, args...)
}

// Debugln log `debug` event with fmt.Sprintln
func (session Session) Debugln(args ...interface{}) {
	session.EventSkip(2).Logln(debugLevel, args...)
}

// Info log `info` event with fmt.Sprint
func (session Session) Info(args ...interface{}) {
	session.EventSkip(2).Log(infoLevel, args...)
}

// Infof log `info` event with fmt.Sprintf
func (session Session) Infof(format string, args ...interface{}) {
	session.EventSkip(2).Logf(infoLevel, format, args...)
}

// Infoln log `info` event with fmt.Sprintln
func (session Session) Infoln(args ...interface{}) {
	session.EventSkip(2).Logln(infoLevel, args...)
}

// Warn log `warn` event with fmt.Sprint
func (session Session) Warn(args ...interface{}) {
	session.EventSkip(2).Log(warnLevel, args...)
}

// Warnf log `warn` event with fmt.Sprintf
func (session Session) Warnf(format string, args ...interface{}) {
	session.EventSkip(2).Logf(warnLevel, format, args...)
}

// Warnln log `warn` event with fmt.Sprintln
func (session Session) Warnln(args ...interface{}) {
	session.EventSkip(2).Logln(warnLevel, args...)
}

// Error log `error` event with fmt.Sprint
func (session Session) Error(args ...interface{}) {
	session.EventSkip(2).Log(errorLevel, args...)
}

// Errorf log `error` event with fmt.Sprintf
func (session Session) Errorf(format string, args ...interface{}) {
	session.EventSkip(2).Logf(errorLevel, format, args...)
}

// Errorln log `error` event with fmt.Sprintln
func (session Session) Errorln(args ...interface{}) {
	session.EventSkip(2).Logln(errorLevel, args...)
}

// Fatal log `fatal` event with fmt.Sprint
func (session Session) Fatal(args ...interface{}) {
	session.EventSkip(2).Log(fatalLevel, args...)
}

// Fatalf log `fatal` event with fmt.Sprintf
func (session Session) Fatalf(format string, args ...interface{}) {
	session.EventSkip(2).Logf(fatalLevel, format, args...)
}

// Fatalln log `fatal` event with fmt.Sprintln
func (session Session) Fatalln(args ...interface{}) {
	session.EventSkip(2).Logln(fatalLevel, args...)
}

// Panic throw `panic` event with fmt.Sprint
func (session Session) Panic(args ...interface{}) {
	event := session.EventSkip(2)
	event.Level = "panic"
	event.Message = fmt.Sprint(args...)
	event.WithField(stackKey, string(debug.Stack()))
	panic(event)
}

// Panicf throw `panic` event with fmt.Sprintf
func (session Session) Panicf(format string, args ...interface{}) {
	event := session.EventSkip(2)
	event.Level = "panic"
	event.Message = fmt.Sprintf(format, args...)
	event.WithField(stackKey, string(debug.Stack()))
	panic(event)
}

// Panicln throw `panic` event with fmt.Sprintln
func (session Session) Panicln(args ...interface{}) {
	event := session.EventSkip(2)
	event.Level = "panic"
	event.Message = fmt.Sprintln(args...)
	event.WithField(stackKey, string(debug.Stack()))
	panic(event)
}

// Log write event with customized level and fmt.Sprint
func (session Session) Log(level string, args ...interface{}) {
	session.EventSkip(2).Log(level, args...)
}

// Logf write event with customized level and fmt.Sprintf
func (session Session) Logf(level string, format string, args ...interface{}) {
	session.EventSkip(2).Logf(level, format, args...)
}

// Logln write event with customized level and fmt.Sprintln
func (session Session) Logln(level string, args ...interface{}) {
	session.EventSkip(2).Logln(level, args...)
}
