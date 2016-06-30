package slog

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestSetErrorKey(t *testing.T) {
	defer func() { errorKey = warnLevel }()
	SetErrorKey("rorre")
	if errorKey != "rorre" {
		t.Errorf("unexpected error key: %q\n", errorKey)
	}
}

func TestSetEmptyErrorKey(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Error("unexpected success")
		}
		errorKey = warnLevel
	}()
	SetErrorKey("")
}

func TestSetStackKey(t *testing.T) {
	defer func() { stackKey = "stack" }()
	SetStackKey("kcats")
	if stackKey != "kcats" {
		t.Errorf("unexpected stack key: %q\n", stackKey)
	}
}

func TestSetEmptyStackKey(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Error("unexpected success")
		}
		stackKey = "stack"
	}()
	SetStackKey("")
}

func TestEventWithError(t *testing.T) {
	event := newEvent(1, nil)
	res := event.WithError(errors.New("test"))
	if res != event {
		t.Error("unexpected return value")
	} else if err := event.Fields[errorKey]; err == nil {
		t.Error("with error fail")
	} else if err, ok := err.(string); !ok {
		t.Error("wrong error type")
	} else if err != "test" {
		t.Error("wrong error value")
	}
}

func TestEventWithField(t *testing.T) {
	event := newEvent(1, nil)
	res := event.WithField("foo", "bar")
	if res != event {
		t.Error("unexpected return value")
	} else if value := event.Fields["foo"]; value == nil {
		t.Error("with field fail")
	} else if value.(string) != "bar" {
		t.Error("wrong value")
	}
}

func TestEventWithFields(t *testing.T) {
	event := newEvent(1, nil)
	res := event.WithFields(Fields{"foo": "bar", "less": "more"})
	if res != event {
		t.Error("unexpected return value")
	} else if value, ok := event.Fields["foo"].(string); !ok || value != "bar" {
		t.Error("wrong value of `foo`:", value)
	} else if value, ok := event.Fields["less"].(string); !ok || value != "more" {
		t.Error("wrong value of `less`:", value)
	}
}

func TestEventFieldify(t *testing.T) {
	timestamp := time.Now()
	globalFields = map[string]interface{}{
		"steve": "jobs",
	}
	event := &Event{
		Timestamp: timestamp,
		Level:     debugLevel,
		Message:   "test",
		Session:   Session{"foo": "bar"},
		Fields:    Fields{"less": "more"},
		Caller: Caller{
			Package: "github.com/yangchenxing/go-slog",
			File:    "event_test.go",
			Line:    72,
			Func:    "TestEventFieldify",
		},
	}
	fields := event.Fieldify(time.RFC3339)
	if len(fields) != 7 {
		t.Error("unexpected result:", fields)
		return
	}
	expectedResult := map[string]interface{}{
		"timestamp": timestamp.Format(time.RFC3339),
		"level":     debugLevel,
		"message":   "test",
		"foo":       "bar",
		"less":      "more",
		"steve":     "jobs",
	}
	for key, value := range expectedResult {
		if fields[key] != value {
			t.Error("unexpected result:", fields)
			return
		}
	}
	if caller, found := fields["caller"]; !found {
		t.Error("miss caller")
		return
	} else if caller != event.Caller {
		t.Error("unexpected result:", fields)
		return
	}
}

func TestEventDebug(t *testing.T) {
	handler := new(receiveHandler)
	handlers = map[string][]Handler{
		debugLevel: []Handler{handler},
	}
	event := newEvent(1, nil)
	event.Debug("test", "event")
	if len(handler.events) != 1 {
		t.Error("unexpected events:", handler.events)
		return
	}
	event = handler.events[0]
	if event.Level != debugLevel || event.Message != fmt.Sprint("test", "event") ||
		event.Caller.File != "event_test.go" || event.Caller.Package != "github.com/yangchenxing/go-slog" ||
		event.Caller.Func != "TestEventDebug" {
		t.Error("unexpected event:", event)
		return
	}
}

func TestEventDebugf(t *testing.T) {
	handler := new(receiveHandler)
	handlers = map[string][]Handler{
		debugLevel: []Handler{handler},
	}
	event := newEvent(1, nil)
	event.Debugf("test %s", "event")
	if len(handler.events) != 1 {
		t.Error("unexpected events:", handler.events)
		return
	}
	event = handler.events[0]
	if event.Level != debugLevel || event.Message != fmt.Sprintf("test %s", "event") ||
		event.Caller.File != "event_test.go" || event.Caller.Package != "github.com/yangchenxing/go-slog" ||
		event.Caller.Func != "TestEventDebugf" {
		t.Error("unexpected event:", event)
		return
	}
}

func TestEventDebugln(t *testing.T) {
	handler := new(receiveHandler)
	handlers = map[string][]Handler{
		debugLevel: []Handler{handler},
	}
	event := newEvent(1, nil)
	event.Debugln("test", "event")
	if len(handler.events) != 1 {
		t.Error("unexpected events:", handler.events)
		return
	}
	event = handler.events[0]
	if event.Level != debugLevel || event.Message != fmt.Sprintln("test", "event") ||
		event.Caller.File != "event_test.go" || event.Caller.Package != "github.com/yangchenxing/go-slog" ||
		event.Caller.Func != "TestEventDebugln" {
		t.Error("unexpected event:", event)
		return
	}
}

func TestEventInfo(t *testing.T) {
	handler := new(receiveHandler)
	handlers = map[string][]Handler{
		infoLevel: []Handler{handler},
	}
	event := newEvent(1, nil)
	event.Info("test", "event")
	if len(handler.events) != 1 {
		t.Error("unexpected events:", handler.events)
		return
	}
	event = handler.events[0]
	if event.Level != infoLevel || event.Message != fmt.Sprint("test", "event") ||
		event.Caller.File != "event_test.go" || event.Caller.Package != "github.com/yangchenxing/go-slog" ||
		event.Caller.Func != "TestEventInfo" {
		t.Error("unexpected event:", event)
		return
	}
}

func TestEventInfof(t *testing.T) {
	handler := new(receiveHandler)
	handlers = map[string][]Handler{
		infoLevel: []Handler{handler},
	}
	event := newEvent(1, nil)
	event.Infof("test %s", "event")
	if len(handler.events) != 1 {
		t.Error("unexpected events:", handler.events)
		return
	}
	event = handler.events[0]
	if event.Level != infoLevel || event.Message != fmt.Sprintf("test %s", "event") ||
		event.Caller.File != "event_test.go" || event.Caller.Package != "github.com/yangchenxing/go-slog" ||
		event.Caller.Func != "TestEventInfof" {
		t.Error("unexpected event:", event)
		return
	}
}

func TestEventInfoln(t *testing.T) {
	handler := new(receiveHandler)
	handlers = map[string][]Handler{
		infoLevel: []Handler{handler},
	}
	event := newEvent(1, nil)
	event.Infoln("test", "event")
	if len(handler.events) != 1 {
		t.Error("unexpected events:", handler.events)
		return
	}
	event = handler.events[0]
	if event.Level != infoLevel || event.Message != fmt.Sprintln("test", "event") ||
		event.Caller.File != "event_test.go" || event.Caller.Package != "github.com/yangchenxing/go-slog" ||
		event.Caller.Func != "TestEventInfoln" {
		t.Error("unexpected event:", event)
		return
	}
}

func TestEventWarn(t *testing.T) {
	handler := new(receiveHandler)
	handlers = map[string][]Handler{
		warnLevel: []Handler{handler},
	}
	event := newEvent(1, nil)
	event.Warn("test", "event")
	if len(handler.events) != 1 {
		t.Error("unexpected events:", handler.events)
		return
	}
	event = handler.events[0]
	if event.Level != warnLevel || event.Message != fmt.Sprint("test", "event") ||
		event.Caller.File != "event_test.go" || event.Caller.Package != "github.com/yangchenxing/go-slog" ||
		event.Caller.Func != "TestEventWarn" {
		t.Error("unexpected event:", event)
		return
	}
}

func TestEventWarnf(t *testing.T) {
	handler := new(receiveHandler)
	handlers = map[string][]Handler{
		warnLevel: []Handler{handler},
	}
	event := newEvent(1, nil)
	event.Warnf("test %s", "event")
	if len(handler.events) != 1 {
		t.Error("unexpected events:", handler.events)
		return
	}
	event = handler.events[0]
	if event.Level != warnLevel || event.Message != fmt.Sprintf("test %s", "event") ||
		event.Caller.File != "event_test.go" || event.Caller.Package != "github.com/yangchenxing/go-slog" ||
		event.Caller.Func != "TestEventWarnf" {
		t.Error("unexpected event:", event)
		return
	}
}

func TestEventWarnln(t *testing.T) {
	handler := new(receiveHandler)
	handlers = map[string][]Handler{
		warnLevel: []Handler{handler},
	}
	event := newEvent(1, nil)
	event.Warnln("test", "event")
	if len(handler.events) != 1 {
		t.Error("unexpected events:", handler.events)
		return
	}
	event = handler.events[0]
	if event.Level != warnLevel || event.Message != fmt.Sprintln("test", "event") ||
		event.Caller.File != "event_test.go" || event.Caller.Package != "github.com/yangchenxing/go-slog" ||
		event.Caller.Func != "TestEventWarnln" {
		t.Error("unexpected event:", event)
		return
	}
}

func TestEventError(t *testing.T) {
	handler := new(receiveHandler)
	handlers = map[string][]Handler{
		errorLevel: []Handler{handler},
	}
	event := newEvent(1, nil)
	event.Error("test", "event")
	if len(handler.events) != 1 {
		t.Error("unexpected events:", handler.events)
		return
	}
	event = handler.events[0]
	if event.Level != errorLevel || event.Message != fmt.Sprint("test", "event") ||
		event.Caller.File != "event_test.go" || event.Caller.Package != "github.com/yangchenxing/go-slog" ||
		event.Caller.Func != "TestEventError" {
		t.Error("unexpected event:", event)
		return
	}
}

func TestEventErrorf(t *testing.T) {
	handler := new(receiveHandler)
	handlers = map[string][]Handler{
		errorLevel: []Handler{handler},
	}
	event := newEvent(1, nil)
	event.Errorf("test %s", "event")
	if len(handler.events) != 1 {
		t.Error("unexpected events:", handler.events)
		return
	}
	event = handler.events[0]
	if event.Level != errorLevel || event.Message != fmt.Sprintf("test %s", "event") ||
		event.Caller.File != "event_test.go" || event.Caller.Package != "github.com/yangchenxing/go-slog" ||
		event.Caller.Func != "TestEventErrorf" {
		t.Error("unexpected event:", event)
		return
	}
}

func TestEventErrorln(t *testing.T) {
	handler := new(receiveHandler)
	handlers = map[string][]Handler{
		errorLevel: []Handler{handler},
	}
	event := newEvent(1, nil)
	event.Errorln("test", "event")
	if len(handler.events) != 1 {
		t.Error("unexpected events:", handler.events)
		return
	}
	event = handler.events[0]
	if event.Level != errorLevel || event.Message != fmt.Sprintln("test", "event") ||
		event.Caller.File != "event_test.go" || event.Caller.Package != "github.com/yangchenxing/go-slog" ||
		event.Caller.Func != "TestEventErrorln" {
		t.Error("unexpected event:", event)
		return
	}
}

func TestEventFatal(t *testing.T) {
	handler := new(receiveHandler)
	handlers = map[string][]Handler{
		"fatal": []Handler{handler},
	}
	event := newEvent(1, nil)
	event.Fatal("test", "event")
	if len(handler.events) != 1 {
		t.Error("unexpected events:", handler.events)
		return
	}
	event = handler.events[0]
	if event.Level != "fatal" || event.Message != fmt.Sprint("test", "event") ||
		event.Caller.File != "event_test.go" || event.Caller.Package != "github.com/yangchenxing/go-slog" ||
		event.Caller.Func != "TestEventFatal" {
		t.Error("unexpected event:", event)
		return
	}
}

func TestEventFatalf(t *testing.T) {
	handler := new(receiveHandler)
	handlers = map[string][]Handler{
		"fatal": []Handler{handler},
	}
	event := newEvent(1, nil)
	event.Fatalf("test %s", "event")
	if len(handler.events) != 1 {
		t.Error("unexpected events:", handler.events)
		return
	}
	event = handler.events[0]
	if event.Level != "fatal" || event.Message != fmt.Sprintf("test %s", "event") ||
		event.Caller.File != "event_test.go" || event.Caller.Package != "github.com/yangchenxing/go-slog" ||
		event.Caller.Func != "TestEventFatalf" {
		t.Error("unexpected event:", event)
		return
	}
}

func TestEventFatalln(t *testing.T) {
	handler := new(receiveHandler)
	handlers = map[string][]Handler{
		"fatal": []Handler{handler},
	}
	event := newEvent(1, nil)
	event.Fatalln("test", "event")
	if len(handler.events) != 1 {
		t.Error("unexpected events:", handler.events)
		return
	}
	event = handler.events[0]
	if event.Level != "fatal" || event.Message != fmt.Sprintln("test", "event") ||
		event.Caller.File != "event_test.go" || event.Caller.Package != "github.com/yangchenxing/go-slog" ||
		event.Caller.Func != "TestEventFatalln" {
		t.Error("unexpected event:", event)
		return
	}
}

func TestEventPanic(t *testing.T) {
	defer func() {
		if event := recover(); event == nil {
			t.Error("no panic")
		} else if event, ok := event.(*Event); !ok {
			t.Error("invalid panic type:", event)
		} else if event.Level != "panic" || event.Message != fmt.Sprint("test", "event") {
			t.Error("unexpected event:", event)
		} else if _, found := event.Fields[stackKey]; !found {
			t.Error("miss stack")
		}
	}()
	event := newEvent(1, nil)
	event.Panic("test", "event")
}

func TestEventPanicf(t *testing.T) {
	defer func() {
		if event := recover(); event == nil {
			t.Error("no panic")
		} else if event, ok := event.(*Event); !ok {
			t.Error("invalid panic type:", event)
		} else if event.Level != "panic" || event.Message != fmt.Sprintf("test %s", "event") {
			t.Error("unexpected event:", event)
		} else if _, found := event.Fields[stackKey]; !found {
			t.Error("miss stack")
		}
	}()
	event := newEvent(1, nil)
	event.Panicf("test %s", "event")
}

func TestEventPanicln(t *testing.T) {
	defer func() {
		if event := recover(); event == nil {
			t.Error("no panic")
		} else if event, ok := event.(*Event); !ok {
			t.Error("invalid panic type:", event)
		} else if event.Level != "panic" || event.Message != fmt.Sprintln("test", "event") {
			t.Error("unexpected event:", event)
		} else if _, found := event.Fields[stackKey]; !found {
			t.Error("miss stack")
		}
	}()
	event := newEvent(1, nil)
	event.Panicln("test", "event")
}

func TestEventDefaultHandler(t *testing.T) {
	defaultHandlerBackup := defaultHandlers
	defer func() { defaultHandlers = defaultHandlerBackup }()
	handler := new(receiveHandler)
	defaultHandlers = map[string][]Handler{
		debugLevel: []Handler{handler},
	}
	handlers = nil
	event := newEvent(1, nil)
	event.Debug("test", "event")
	if len(handler.events) != 1 {
		t.Error("unexpected events:", handler.events)
		return
	}
	event = handler.events[0]
	if event.Level != debugLevel || event.Message != fmt.Sprint("test", "event") ||
		event.Caller.File != "event_test.go" || event.Caller.Package != "github.com/yangchenxing/go-slog" ||
		event.Caller.Func != "TestEventDefaultHandler" {
		t.Error("unexpected event:", event)
		return
	}
}
