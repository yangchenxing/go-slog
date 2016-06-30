package slog

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/smtp"
	"os"
	"strings"
	"sync"
	"time"
)

type PlainTextEmailHandler struct {
	sync.Mutex
	SMTPServer            string
	SMTPUsername          string
	SMTPPassword          string
	Sender                string
	Receivers             []string
	Subject               string
	ContentFormatter      MultiEventFormatter
	AggregationTimeWindow time.Duration
	eventChan             chan *Event
}

func (handler *PlainTextEmailHandler) Handle(event *Event) {
	handler.Lock()
	defer handler.Unlock()
	if handler.eventChan == nil {
		handler.eventChan = make(chan *Event, 16)
		go handler.aggregateEvents()
	}
	handler.eventChan <- event
}

func (handler *PlainTextEmailHandler) aggregateEvents() {
	events := make([]*Event, 0, 16)
	for {
		events = events[:0]
		event := <-handler.eventChan
		events = append(events, event)
		ticker := time.NewTicker(handler.AggregationTimeWindow)
	AggregationLoop:
		for {
			select {
			case event := <-handler.eventChan:
				events = append(events, event)
			case <-ticker.C:
				ticker.Stop()
				break AggregationLoop
			}
		}
		var buffer bytes.Buffer
		fmt.Fprintf(&buffer, "From: %s\r\n", handler.Sender)
		fmt.Fprintf(&buffer, "To: %s\r\n", strings.Join(handler.Receivers, ","))
		fmt.Fprintf(&buffer, "Subject: %s\r\n",
			"=?utf-8?B?"+base64.StdEncoding.EncodeToString([]byte(handler.Subject))+"?=")
		buffer.WriteString("Content-Type: text/plain; charset=UTF-8\r\n\r\n")
		body, err := handler.ContentFormatter.FormatEvents(events)
		if err != nil {
			fmt.Fprintln(os.Stderr, "format email body fail:", err.Error())
			continue
		}
		buffer.Write(body)
		smtpHost := strings.Split(handler.SMTPServer, ":")[0]
		auth := smtp.PlainAuth("", handler.SMTPUsername, handler.SMTPPassword, smtpHost)
		if err := smtp.SendMail(handler.SMTPServer, auth, handler.Sender, handler.Receivers, buffer.Bytes()); err != nil {
			fmt.Fprintln(os.Stderr, "send mail fail:", err.Error())
		}
	}
}
