package slog

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

type receiveHandler struct {
	events []*Event
}

func (handler *receiveHandler) Handle(event *Event) {
	if handler.events == nil {
		handler.events = make([]*Event, 0, 1)
	}
	handler.events = append(handler.events, event)
}

func TestAddHandler(t *testing.T) {
	handlers = nil
	handler1 := new(receiveHandler)
	handler2 := new(receiveHandler)
	AddHandler([]string{"info", "debug"}, handler1)
	AddHandler([]string{"debug", "info"}, handler2)
	infoHandlers := handlers["info"]
	debugHandlers := handlers["debug"]
	if len(infoHandlers) != 2 || infoHandlers[0] != handler1 || infoHandlers[1] != handler2 ||
		len(debugHandlers) != 2 || debugHandlers[0] != handler1 || debugHandlers[1] != handler2 {
		t.Error("unexpected handlers:", handlers)
	}
}

func TestJsonHandler(t *testing.T) {
	event := newEvent(1, nil)
	event.Level = "debug"
	event.Message = "test"
	writer := new(bufferWriter)
	handler := &JsonHandler{
		Writer: writer,
	}
	handler.Handle(event)
	var eventJSON map[string]interface{}
	if err := json.Unmarshal(writer.Bytes(), &eventJSON); err != nil {
		t.Error("unmarshal json fail:", err.Error())
	} else if eventJSON["level"] != "debug" || eventJSON["message"] != "test" {
		t.Error("unmarshal json fail:", err.Error())
	}
}

func TestJsonHandlerFail1(t *testing.T) {
	os.Remove("temp.txt")
	tempfile, err := os.OpenFile("temp.txt", os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		t.Error("open file fail:", err.Error())
		return
	}
	defer func() {
		tempfile.Close()
		os.Remove("temp.txt")
	}()
	stderr := os.Stderr
	os.Stderr = tempfile
	defer func() {
		os.Stderr = stderr
	}()
	writer := new(errorWriter)
	writer.err = errors.New("_error_")
	event := newEvent(1, nil).WithField("ch", make(chan bool))
	handler := &JsonHandler{
		Writer: writer,
	}
	handler.Handle(event)
	tempfile.Close()
	content, err := ioutil.ReadFile("temp.txt")
	if err != nil {
		t.Error("read temp.txt fail:", err.Error())
		return
	} else if !strings.HasPrefix(string(content), "marshal json fail:") {
		t.Error("unexpected error message:", string(content))
		return
	}
}

func TestJsonHandlerFail2(t *testing.T) {
	os.Remove("temp.txt")
	tempfile, err := os.OpenFile("temp.txt", os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		t.Error("open file fail:", err.Error())
		return
	}
	defer func() {
		tempfile.Close()
		os.Remove("temp.txt")
	}()
	stderr := os.Stderr
	os.Stderr = tempfile
	defer func() {
		os.Stderr = stderr
	}()
	writer := new(errorWriter)
	writer.err = errors.New("_error_")
	event := newEvent(1, nil)
	handler := &JsonHandler{
		Writer: writer,
	}
	handler.Handle(event)
	tempfile.Close()
	content, err := ioutil.ReadFile("temp.txt")
	if err != nil {
		t.Error("read temp.txt fail:", err.Error())
		return
	} else if !strings.HasPrefix(string(content), "write json fail:") {
		t.Error("unexpected error message:", string(content))
		return
	}
}

func TestPlainTextHandler(t *testing.T) {
	event := newEvent(1, nil)
	event.Level = "debug"
	event.Message = "test"
	writer := new(bufferWriter)
	handler := &PlainTextHandler{
		Formatter: &PlainTextFormatter{
			EventFormat: "%(level|s) %(caller.package|s)/%(caller.file|s):%(caller.func|s) %(message|s)",
		},
		Writer: writer,
	}
	handler.Handle(event)
	if writer.String() != "debug github.com/yangchenxing/go-slog/handler_test.go:TestPlainTextHandler test" {
		t.Error("unexpected result:", writer.String())
	}
}

func TestPlainTextHandlerFail(t *testing.T) {
	os.Remove("temp.txt")
	tempfile, err := os.OpenFile("temp.txt", os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		t.Error("open file fail:", err.Error())
		return
	}
	defer func() {
		tempfile.Close()
		os.Remove("temp.txt")
	}()
	stderr := os.Stderr
	os.Stderr = tempfile
	defer func() {
		os.Stderr = stderr
	}()
	writer := new(errorWriter)
	writer.err = errors.New("_error_")
	event := newEvent(1, nil)
	handler := &PlainTextHandler{
		Formatter: &PlainTextFormatter{
			EventFormat: "%(level|s) %(caller.package|s)/%(caller.file|s):%(caller.func|s) %(message|s)",
		},
		Writer: writer,
	}
	handler.Handle(event)
	tempfile.Close()
	content, err := ioutil.ReadFile("temp.txt")
	if err != nil {
		t.Error("read temp.txt fail:", err.Error())
		return
	} else if !strings.HasPrefix(string(content), "write text fail:") {
		t.Error("unexpected error message:", string(content))
		return
	}
}
