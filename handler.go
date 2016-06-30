package slog

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

var (
	handlers       map[string][]Handler
	defaultHandler = &PlainTextHandler{
		Formatter: &PlainTextFormatter{
			EventFormat: "%(level|s) [%(timestamp|s)] %(message|s) [%(.all_fields_space_seperated_text|s)]",
		},
		Writer: StdoutWriter,
	}
	defaultHandlers = map[string][]Handler{
		"debug": []Handler{defaultHandler},
		"info":  []Handler{defaultHandler},
		"warn":  []Handler{defaultHandler},
		"error": []Handler{defaultHandler},
		"fatal": []Handler{defaultHandler},
	}
)

type Handler interface {
	Handle(*Event)
}

func AddHandler(levels []string, handler Handler) {
	if handlers == nil {
		handlers = make(map[string][]Handler)
	}
	for _, level := range levels {
		if levelHandlers := handlers[level]; levelHandlers != nil {
			handlers[level] = append(levelHandlers, handler)
		} else {
			handlers[level] = []Handler{handler}
		}
	}
}

type JsonHandler struct {
	TimestampFormat string
	Writer          Writer
}

func (handler *JsonHandler) Handle(event *Event) {
	timestampFormat := handler.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.RFC3339
	}
	if content, err := json.Marshal(event.Fieldify(handler.TimestampFormat)); err != nil {
		fmt.Fprintf(os.Stderr, "marshal json fail: event=%v, error=%q\n", event, err.Error())
	} else if err := handler.Writer.Write(content); err != nil {
		fmt.Fprintf(os.Stderr, "write json fail: event=%v, error=%q\n", event, err.Error())
	}
}

type PlainTextHandler struct {
	Formatter *PlainTextFormatter
	Writer    Writer
}

func (handler *PlainTextHandler) Handle(event *Event) {
	content, _ := handler.Formatter.FormatEvent(event)
	if err := handler.Writer.Write(content); err != nil {
		fmt.Fprintf(os.Stderr, "write text fail: event=%v, error=%q\n", event, err.Error())
	}
}
