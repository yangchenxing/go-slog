package slog

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

type bufferWriter struct {
	bytes.Buffer
}

func (writer *bufferWriter) Write(b []byte) error {
	writer.Buffer.Write(b)
	return nil
}

type errorWriter struct {
	err error
}

func (writer *errorWriter) Write([]byte) error {
	return writer.err
}

func TestFileWriter(t *testing.T) {
	os.Remove("temp.txt")
	tempfile, err := os.OpenFile("temp.txt", os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		t.Error("create tempfile fail:", err.Error())
	}
	defer func() {
		os.Remove("temp.txt")
	}()
	writer := &FileWriter{File: tempfile}
	writer.Write([]byte("test"))
	tempfile.Close()
	content, err := ioutil.ReadFile("temp.txt")
	if err != nil {
		t.Error("read tempfile fail:", err.Error())
	} else if string(content) != "test\n" {
		t.Error("unexpected file content:", string(content))
	}
}
