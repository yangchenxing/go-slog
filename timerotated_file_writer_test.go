package slog

import (
	"os"
	"testing"
	"time"
)

func TestTimeRotatedFileWriterRotate(t *testing.T) {
	writer := &TimeRotatedFileWriter{
		Path:            "temp.log",
		TimestampFormat: "20060102150405",
		Interval:        time.Second,
		Keep:            time.Minute,
	}
	if err := writer.open(); err != nil {
		t.Error("open fail:", err.Error())
		return
	}
	defer func() {
		os.Remove("temp.log")
	}()
	timestamp := time.Now().Truncate(time.Second).Add(-1 * time.Second)
	writer.rotate(timestamp)
}
