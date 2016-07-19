package slog

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	tzOffset     time.Duration
	tzOffsetBack time.Duration
	newline      = []byte("\n")
)

func init() {
	_, offset := time.Now().Zone()
	tzOffset = time.Second * time.Duration(offset)
	tzOffsetBack = -1 * tzOffset
}

type TimeRotatedFileWriter struct {
	sync.Mutex
	Path            string
	TimestampFormat string
	Interval        time.Duration
	Keep            time.Duration
	file            *os.File
	serving         bool
	servingStop     chan time.Time
	servingStopped  chan time.Time
}

func (writer *TimeRotatedFileWriter) Write(content []byte) error {
	if writer.file == nil {
		if err := writer.open(); err != nil {
			return fmt.Errorf("open file fail: %s", err.Error())
		}
	}
	if !writer.serving {
		writer.startMaid()
	}
	writer.file.Write(content)
	_, err := writer.file.Write(newline)
	return err
}

func (writer *TimeRotatedFileWriter) open() error {
	writer.Lock()
	defer writer.Unlock()
	if writer.file != nil {
		return nil
	}
	var err error
	writer.file, err = os.OpenFile(writer.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	return err
}

func (writer *TimeRotatedFileWriter) startMaid() {
	writer.Lock()
	defer writer.Unlock()
	if writer.serving {
		return
	}
	go writer.serve()
	writer.serving = true
}

func (writer *TimeRotatedFileWriter) initializeChans() {
	writer.Lock()
	defer writer.Unlock()
	if writer.servingStop == nil {
		writer.servingStop = make(chan time.Time)
	}
	if writer.servingStopped == nil {
		writer.servingStopped = make(chan time.Time)
	}
}

func (writer *TimeRotatedFileWriter) stopServe() {
	writer.initializeChans()
	writer.servingStop <- time.Now()
}

func (writer *TimeRotatedFileWriter) serve() {
	writer.initializeChans()
	currentTimestamp := time.Now().Add(tzOffsetBack).Truncate(writer.Interval).Add(tzOffset)
ForLoop:
	for {
		nextTimestamp := currentTimestamp.Add(writer.Interval)
		select {
		case <-time.After(nextTimestamp.Sub(time.Now())):
			// 切割和清理同周期进行
			if err := writer.rotate(currentTimestamp); err != nil {
				fmt.Fprintf(os.Stderr, "rotate log file %q fail: %s\n", writer.Path, err.Error())
			}
			if err := writer.clean(currentTimestamp); err != nil {
				fmt.Fprintf(os.Stderr, "clean log file %q fail: %s\n", writer.Path, err.Error())
			}
			currentTimestamp = nextTimestamp
		case <-writer.servingStop:
			break ForLoop
		}
	}
}

func (writer *TimeRotatedFileWriter) rotate(timestamp time.Time) error {
	writer.Lock()
	defer writer.Unlock()
	splitPath := writer.Path + "." + timestamp.Format(writer.TimestampFormat)
	if err := os.Rename(writer.Path, splitPath); err != nil {
		return fmt.Errorf("rename %q to %q fail: %s", writer.Path, splitPath, err.Error())
	}
	oldFile := writer.file
	writer.file = nil
	// 延迟1秒关闭上一个文件，避免对写日志操作加锁
	go func() {
		time.Sleep(time.Second)
		oldFile.Close()
	}()
	return nil
}

func (writer *TimeRotatedFileWriter) clean(timestamp time.Time) error {
	// 未指定保存时长时不清理
	if writer.Keep <= 0 {
		return nil
	}
	writer.Lock()
	defer writer.Unlock()
	expireTimestamp := timestamp.Add(-1 * writer.Keep)
	dir := filepath.Dir(writer.Path)
	infos, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read log dir %q fail: %s", dir, err.Error())
	}
	prefix := filepath.Base(writer.Path) + "."
	for _, info := range infos {
		if strings.HasPrefix(info.Name(), prefix) {
			splitTimestamp, err := time.Parse(writer.TimestampFormat, info.Name()[len(prefix):])
			if err != nil {
				return fmt.Errorf("parse timestamp of %q fail: %s",
					filepath.Join(dir, info.Name()), err.Error())
			}
			if !splitTimestamp.After(expireTimestamp) {
				splitPath := filepath.Join(dir, info.Name())
				if err := os.Remove(splitPath); err != nil {
					return fmt.Errorf("remove log file %q fail: %s", splitPath, err.Error())
				}
			}
		}
	}
	return nil
}
