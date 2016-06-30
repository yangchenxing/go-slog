package slog

import (
	"os"
)

type Writer interface {
	Write([]byte) error
}

type FileWriter struct {
	File *os.File
}

func (writer FileWriter) Write(content []byte) error {
	_, err := writer.File.Write(content)
	writer.File.WriteString("\n")
	if err == nil {
		writer.File.Sync()
	}
	return err
}

var (
	StdoutWriter = &FileWriter{File: os.Stdout}
	StderrWriter = &FileWriter{File: os.Stderr}
)
