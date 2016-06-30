package slog

import (
	"fmt"
	"testing"
)

func TestSessionWithField(t *testing.T) {
	session := NewSession()
	session.WithField("foo", "bar")
	if len(session) != 1 || session["foo"] != "bar" {
		t.Error("unexpected result:", session)
	}
}

func TestSessionWithFields(t *testing.T) {
	session := NewSession()
	session.WithFields(Fields{"foo": "bar", "less": "more"})
	if len(session) != 2 || session["foo"] != "bar" || session["less"] != "more" {
		t.Error("unexpected result:", session)
	}
}

func TestSessionNewEvent(t *testing.T) {
	session := NewSession()
	event := session.Event()
	caller := event.Caller
	if caller.Func != "TestSessionNewEvent" || caller.Package != "github.com/yangchenxing/go-slog" ||
		caller.File != "session_test.go" {
		t.Error("unexpected result:", session)
	}
}

func TestSessionDebug(t *testing.T) {
	handler := new(receiveHandler)
	handlers = map[string][]Handler{
		debugLevel: []Handler{handler},
	}
	session := NewSession()
	session.Debug("test", "event")
	if len(handler.events) != 1 {
		t.Error("unexpected events:", handler.events)
		return
	}
	event := handler.events[0]
	if event.Level != debugLevel || event.Message != fmt.Sprint("test", "event") ||
		event.Caller.File != "session_test.go" || event.Caller.Package != "github.com/yangchenxing/go-slog" ||
		event.Caller.Func != "TestSessionDebug" {
		t.Error("unexpected event:", event)
		return
	}
}

func TestSessionDebugf(t *testing.T) {
	handler := new(receiveHandler)
	handlers = map[string][]Handler{
		debugLevel: []Handler{handler},
	}
	session := NewSession()
	session.Debugf("test %s", "event")
	if len(handler.events) != 1 {
		t.Error("unexpected events:", handler.events)
		return
	}
	event := handler.events[0]
	if event.Level != debugLevel || event.Message != fmt.Sprintf("test %s", "event") ||
		event.Caller.File != "session_test.go" || event.Caller.Package != "github.com/yangchenxing/go-slog" ||
		event.Caller.Func != "TestSessionDebugf" {
		t.Error("unexpected event:", event)
		return
	}
}

func TestSessionDebugln(t *testing.T) {
	handler := new(receiveHandler)
	handlers = map[string][]Handler{
		debugLevel: []Handler{handler},
	}
	session := NewSession()
	session.Debugln("test", "event")
	if len(handler.events) != 1 {
		t.Error("unexpected events:", handler.events)
		return
	}
	event := handler.events[0]
	if event.Level != debugLevel || event.Message != fmt.Sprintln("test", "event") ||
		event.Caller.File != "session_test.go" || event.Caller.Package != "github.com/yangchenxing/go-slog" ||
		event.Caller.Func != "TestSessionDebugln" {
		t.Error("unexpected event:", event)
		return
	}
}

func TestSessionInfo(t *testing.T) {
	handler := new(receiveHandler)
	handlers = map[string][]Handler{
		infoLevel: []Handler{handler},
	}
	session := NewSession()
	session.Info("test", "event")
	if len(handler.events) != 1 {
		t.Error("unexpected events:", handler.events)
		return
	}
	event := handler.events[0]
	if event.Level != infoLevel || event.Message != fmt.Sprint("test", "event") ||
		event.Caller.File != "session_test.go" || event.Caller.Package != "github.com/yangchenxing/go-slog" ||
		event.Caller.Func != "TestSessionInfo" {
		t.Error("unexpected event:", event)
		return
	}
}

func TestSessionInfof(t *testing.T) {
	handler := new(receiveHandler)
	handlers = map[string][]Handler{
		infoLevel: []Handler{handler},
	}
	session := NewSession()
	session.Infof("test %s", "event")
	if len(handler.events) != 1 {
		t.Error("unexpected events:", handler.events)
		return
	}
	event := handler.events[0]
	if event.Level != infoLevel || event.Message != fmt.Sprintf("test %s", "event") ||
		event.Caller.File != "session_test.go" || event.Caller.Package != "github.com/yangchenxing/go-slog" ||
		event.Caller.Func != "TestSessionInfof" {
		t.Error("unexpected event:", event)
		return
	}
}

func TestSessionInfoln(t *testing.T) {
	handler := new(receiveHandler)
	handlers = map[string][]Handler{
		infoLevel: []Handler{handler},
	}
	session := NewSession()
	session.Infoln("test", "event")
	if len(handler.events) != 1 {
		t.Error("unexpected events:", handler.events)
		return
	}
	event := handler.events[0]
	if event.Level != infoLevel || event.Message != fmt.Sprintln("test", "event") ||
		event.Caller.File != "session_test.go" || event.Caller.Package != "github.com/yangchenxing/go-slog" ||
		event.Caller.Func != "TestSessionInfoln" {
		t.Error("unexpected event:", event)
		return
	}
}

func TestSessionWarn(t *testing.T) {
	handler := new(receiveHandler)
	handlers = map[string][]Handler{
		warnLevel: []Handler{handler},
	}
	session := NewSession()
	session.Warn("test", "event")
	if len(handler.events) != 1 {
		t.Error("unexpected events:", handler.events)
		return
	}
	event := handler.events[0]
	if event.Level != warnLevel || event.Message != fmt.Sprint("test", "event") ||
		event.Caller.File != "session_test.go" || event.Caller.Package != "github.com/yangchenxing/go-slog" ||
		event.Caller.Func != "TestSessionWarn" {
		t.Error("unexpected event:", event)
		return
	}
}

func TestSessionWarnf(t *testing.T) {
	handler := new(receiveHandler)
	handlers = map[string][]Handler{
		warnLevel: []Handler{handler},
	}
	session := NewSession()
	session.Warnf("test %s", "event")
	if len(handler.events) != 1 {
		t.Error("unexpected events:", handler.events)
		return
	}
	event := handler.events[0]
	if event.Level != warnLevel || event.Message != fmt.Sprintf("test %s", "event") ||
		event.Caller.File != "session_test.go" || event.Caller.Package != "github.com/yangchenxing/go-slog" ||
		event.Caller.Func != "TestSessionWarnf" {
		t.Error("unexpected event:", event)
		return
	}
}

func TestSessionWarnln(t *testing.T) {
	handler := new(receiveHandler)
	handlers = map[string][]Handler{
		warnLevel: []Handler{handler},
	}
	session := NewSession()
	session.Warnln("test", "event")
	if len(handler.events) != 1 {
		t.Error("unexpected events:", handler.events)
		return
	}
	event := handler.events[0]
	if event.Level != warnLevel || event.Message != fmt.Sprintln("test", "event") ||
		event.Caller.File != "session_test.go" || event.Caller.Package != "github.com/yangchenxing/go-slog" ||
		event.Caller.Func != "TestSessionWarnln" {
		t.Error("unexpected event:", event)
		return
	}
}

func TestSessionError(t *testing.T) {
	handler := new(receiveHandler)
	handlers = map[string][]Handler{
		errorLevel: []Handler{handler},
	}
	session := NewSession()
	session.Error("test", "event")
	if len(handler.events) != 1 {
		t.Error("unexpected events:", handler.events)
		return
	}
	event := handler.events[0]
	if event.Level != errorLevel || event.Message != fmt.Sprint("test", "event") ||
		event.Caller.File != "session_test.go" || event.Caller.Package != "github.com/yangchenxing/go-slog" ||
		event.Caller.Func != "TestSessionError" {
		t.Error("unexpected event:", event)
		return
	}
}

func TestSessionErrorf(t *testing.T) {
	handler := new(receiveHandler)
	handlers = map[string][]Handler{
		errorLevel: []Handler{handler},
	}
	session := NewSession()
	session.Errorf("test %s", "event")
	if len(handler.events) != 1 {
		t.Error("unexpected events:", handler.events)
		return
	}
	event := handler.events[0]
	if event.Level != errorLevel || event.Message != fmt.Sprintf("test %s", "event") ||
		event.Caller.File != "session_test.go" || event.Caller.Package != "github.com/yangchenxing/go-slog" ||
		event.Caller.Func != "TestSessionErrorf" {
		t.Error("unexpected event:", event)
		return
	}
}

func TestSessionErrorln(t *testing.T) {
	handler := new(receiveHandler)
	handlers = map[string][]Handler{
		errorLevel: []Handler{handler},
	}
	session := NewSession()
	session.Errorln("test", "event")
	if len(handler.events) != 1 {
		t.Error("unexpected events:", handler.events)
		return
	}
	event := handler.events[0]
	if event.Level != errorLevel || event.Message != fmt.Sprintln("test", "event") ||
		event.Caller.File != "session_test.go" || event.Caller.Package != "github.com/yangchenxing/go-slog" ||
		event.Caller.Func != "TestSessionErrorln" {
		t.Error("unexpected event:", event)
		return
	}
}

func TestSessionFatal(t *testing.T) {
	handler := new(receiveHandler)
	handlers = map[string][]Handler{
		fatalLevel: []Handler{handler},
	}
	session := NewSession()
	session.Fatal("test", "event")
	if len(handler.events) != 1 {
		t.Error("unexpected events:", handler.events)
		return
	}
	event := handler.events[0]
	if event.Level != fatalLevel || event.Message != fmt.Sprint("test", "event") ||
		event.Caller.File != "session_test.go" || event.Caller.Package != "github.com/yangchenxing/go-slog" ||
		event.Caller.Func != "TestSessionFatal" {
		t.Error("unexpected event:", event)
		return
	}
}

func TestSessionFatalf(t *testing.T) {
	handler := new(receiveHandler)
	handlers = map[string][]Handler{
		fatalLevel: []Handler{handler},
	}
	session := NewSession()
	session.Fatalf("test %s", "event")
	if len(handler.events) != 1 {
		t.Error("unexpected events:", handler.events)
		return
	}
	event := handler.events[0]
	if event.Level != fatalLevel || event.Message != fmt.Sprintf("test %s", "event") ||
		event.Caller.File != "session_test.go" || event.Caller.Package != "github.com/yangchenxing/go-slog" ||
		event.Caller.Func != "TestSessionFatalf" {
		t.Error("unexpected event:", event)
		return
	}
}

func TestSessionFatalln(t *testing.T) {
	handler := new(receiveHandler)
	handlers = map[string][]Handler{
		fatalLevel: []Handler{handler},
	}
	session := NewSession()
	session.Fatalln("test", "event")
	if len(handler.events) != 1 {
		t.Error("unexpected events:", handler.events)
		return
	}
	event := handler.events[0]
	if event.Level != fatalLevel || event.Message != fmt.Sprintln("test", "event") ||
		event.Caller.File != "session_test.go" || event.Caller.Package != "github.com/yangchenxing/go-slog" ||
		event.Caller.Func != "TestSessionFatalln" {
		t.Error("unexpected event:", event)
		return
	}
}

func TestSessionPanic(t *testing.T) {
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
	session := NewSession()
	session.Panic("test", "event")
}

func TestSessionPanicf(t *testing.T) {
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
	session := NewSession()
	session.Panicf("test %s", "event")
}

func TestSessionPanicln(t *testing.T) {
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
	session := NewSession()
	session.Panicln("test", "event")
}

func TestSessionLog(t *testing.T) {
	handler := new(receiveHandler)
	handlers = map[string][]Handler{
		"log": []Handler{handler},
	}
	session := NewSession()
	session.Log("log", "test", "event")
	if len(handler.events) != 1 {
		t.Error("unexpected events:", handler.events)
		return
	}
	event := handler.events[0]
	if event.Level != "log" || event.Message != fmt.Sprint("test", "event") ||
		event.Caller.File != "session_test.go" || event.Caller.Package != "github.com/yangchenxing/go-slog" ||
		event.Caller.Func != "TestSessionLog" {
		t.Error("unexpected event:", event)
		return
	}
}

func TestSessionLogf(t *testing.T) {
	handler := new(receiveHandler)
	handlers = map[string][]Handler{
		"log": []Handler{handler},
	}
	session := NewSession()
	session.Logf("log", "test %s", "event")
	if len(handler.events) != 1 {
		t.Error("unexpected events:", handler.events)
		return
	}
	event := handler.events[0]
	if event.Level != "log" || event.Message != fmt.Sprintf("test %s", "event") ||
		event.Caller.File != "session_test.go" || event.Caller.Package != "github.com/yangchenxing/go-slog" ||
		event.Caller.Func != "TestSessionLogf" {
		t.Error("unexpected event:", event)
		return
	}
}

func TestSessionLogln(t *testing.T) {
	handler := new(receiveHandler)
	handlers = map[string][]Handler{
		"log": []Handler{handler},
	}
	session := NewSession()
	session.Logln("log", "test", "event")
	if len(handler.events) != 1 {
		t.Error("unexpected events:", handler.events)
		return
	}
	event := handler.events[0]
	if event.Level != "log" || event.Message != fmt.Sprintln("test", "event") ||
		event.Caller.File != "session_test.go" || event.Caller.Package != "github.com/yangchenxing/go-slog" ||
		event.Caller.Func != "TestSessionLogln" {
		t.Error("unexpected event:", event)
		return
	}
}
