package slog

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"runtime"
	"runtime/debug"
	"time"
)

var (
	debugLevel = "debug"
	infoLevel  = "info"
	warnLevel  = "warn"
	errorLevel = "error"
	fatalLevel = "fatal"
)

// Caller provide caller information of log event
type Caller struct {
	Package string `json:"package"`
	File    string `json:"file"`
	Func    string `json:"func"`
	Line    int    `json:"line"`
}

// Event represents the log event
type Event struct {
	Timestamp time.Time
	Level     string
	Message   string
	Session   Session
	Fields    Fields
	Caller    Caller
}

// Fields type, used by `WithFields`
type Fields map[string]interface{}

var (
	callerCache        = make(map[uintptr]Caller)
	errorKey           = "error"
	stackKey           = "stack"
	funcnamePattern, _ = regexp.Compile("(\\.[^/.]+)+$")
)

// SetErrorKey change field key for error
func SetErrorKey(key string) {
	if key == "" {
		panic(errors.New("error key can not be empty string"))
	}
	errorKey = key
}

// SetStackKey change field key for panic stack
func SetStackKey(key string) {
	if key == "" {
		panic(errors.New("stack key can not be empty string"))
	}
	stackKey = key
}

func newEvent(skip int, session Session) *Event {
	event := &Event{
		Timestamp: time.Now(),
		Fields:    make(map[string]interface{}),
	}
	// 获取Caller信息
	pc, file, line, ok := runtime.Caller(skip)
	if ok {
		caller, found := callerCache[pc]
		if !found {
			if ffp := runtime.FuncForPC(pc); ffp != nil {
				fullfunc := ffp.Name()
				caller.Func = funcnamePattern.FindString(fullfunc)[1:]
				caller.Package = fullfunc[:len(fullfunc)-len(caller.Func)-1]
			}
			caller.File = filepath.Base(file)
			caller.Line = line
			callerCache[pc] = caller
		}
		event.Caller = caller
	}
	event.Session = session
	return event
}

// WithError add err.Error() to event.Fields with errorKey
func (event *Event) WithError(err error) *Event {
	event.Fields[errorKey] = err.Error()
	return event
}

// WithField add value to event.Fields with key
func (event *Event) WithField(key string, value interface{}) *Event {
	event.Fields[key] = value
	return event
}

// WithFields add multiple key value pairs to event.Fields
func (event *Event) WithFields(fields Fields) *Event {
	for key, value := range fields {
		event.Fields[key] = value
	}
	return event
}

// Debug log `debug` event with fmt.Sprint
func (event *Event) Debug(args ...interface{}) {
	event.Log(debugLevel, args...)
}

// Debugf log `debug` event with fmt.Sprintf
func (event *Event) Debugf(format string, args ...interface{}) {
	event.Logf(debugLevel, format, args...)
}

// Debugln log `debug` event with fmt.Sprintln
func (event *Event) Debugln(args ...interface{}) {
	event.Logln(debugLevel, args...)
}

// Info log `info` event with fmt.Sprint
func (event *Event) Info(args ...interface{}) {
	event.Log(infoLevel, args...)
}

// Infof log `info` event with fmt.Sprintf
func (event *Event) Infof(format string, args ...interface{}) {
	event.Logf(infoLevel, format, args...)
}

// Infoln log `info` event with fmt.Sprintln
func (event *Event) Infoln(args ...interface{}) {
	event.Logln(infoLevel, args...)
}

// Warn log `warn` event with fmt.Sprint
func (event *Event) Warn(args ...interface{}) {
	event.Log(warnLevel, args...)
}

// Warnf log `warn` event with fmt.Sprintf
func (event *Event) Warnf(format string, args ...interface{}) {
	event.Logf(warnLevel, format, args...)
}

// Warnln log `warn` event with fmt.Sprintln
func (event *Event) Warnln(args ...interface{}) {
	event.Logln(warnLevel, args...)
}

// Error log `error` event with fmt.Sprint
func (event *Event) Error(args ...interface{}) {
	event.Log(errorLevel, args...)
}

// Errorf log `error` event with fmt.Sprintf
func (event *Event) Errorf(format string, args ...interface{}) {

	event.Logf(errorLevel, format, args...)
}

// Errorln log `error` event with fmt.Sprintln
func (event *Event) Errorln(args ...interface{}) {
	event.Logln(errorLevel, args...)
}

// Fatal log `fatal` event with fmt.Sprint
func (event *Event) Fatal(args ...interface{}) {
	event.Log(fatalLevel, args...)
}

// Fatalf log `fatal` event with fmt.Sprintf
func (event *Event) Fatalf(format string, args ...interface{}) {
	event.Logf(fatalLevel, format, args...)
}

// Fatalln log `fatal` event with fmt.Sprintln
func (event *Event) Fatalln(args ...interface{}) {
	event.Logln(fatalLevel, args...)
}

// Panic throw `panic` event with fmt.Sprint
func (event *Event) Panic(args ...interface{}) {
	event.Level = "panic"
	event.Message = fmt.Sprint(args...)
	event.WithField(stackKey, string(debug.Stack()))
	panic(event)
}

// Panicf throw `panic` event with fmt.Sprintf
func (event *Event) Panicf(format string, args ...interface{}) {
	event.Level = "panic"
	event.Message = fmt.Sprintf(format, args...)
	event.WithField(stackKey, string(debug.Stack()))
	panic(event)
}

// Panicln throw `panic` event with fmt.Sprintln
func (event *Event) Panicln(args ...interface{}) {
	event.Level = "panic"
	event.Message = fmt.Sprintln(args...)
	event.WithField(stackKey, string(debug.Stack()))
	panic(event)
}

// Log write event with customized level and fmt.Sprint
func (event *Event) Log(level string, args ...interface{}) {
	event.Message = fmt.Sprint(args...)
	event.Level = level
	event.write()
}

// Logf write event with customized level and fmt.Sprintf
func (event *Event) Logf(level string, format string, args ...interface{}) {
	event.Message = fmt.Sprintf(format, args...)
	event.Level = level
	event.write()
}

// Logln write event with customized level and fmt.Sprintln
func (event *Event) Logln(level string, args ...interface{}) {
	event.Message = fmt.Sprintln(args...)
	event.Level = level
	event.write()
}

// Fieldify convert event to map[string]interface{}
func (event *Event) Fieldify(timestampFormat string) map[string]interface{} {
	fields := make(map[string]interface{})
	for key, value := range globalFields {
		fields[key] = value
	}
	for key, value := range event.Session {
		fields[key] = value
	}
	for key, value := range event.Fields {
		fields[key] = value
	}
	fields["timestamp"] = time.Now().Format(timestampFormat)
	fields["level"] = event.Level
	fields["message"] = event.Message
	if timestampFormat != "" {
		fields["timestamp"] = event.Timestamp.Format(timestampFormat)
	}
	fields["caller"] = event.Caller
	return fields
}

func (event *Event) write() {
	var levelHandlers []Handler
	if handlers == nil {
		levelHandlers = defaultHandlers[event.Level]
	} else {
		levelHandlers = handlers[event.Level]
	}
	for _, handler := range levelHandlers {
		handler.Handle(event)
	}
}
