package slog

import (
	"strings"
	"testing"
	"time"

	"github.com/fatih/color"
)

func TestJoinFields(t *testing.T) {
	fields := map[string]interface{}{
		"int":  1,
		"str":  "foo",
		"bool": true,
	}
	res := JoinFields(fields, "=", " ", true)
	if res != "bool=true int=1 str=\"foo\"" {
		t.Error("unexpected result:", res)
	}
}

func TestInitializePlainTextFormatter(t *testing.T) {
	var formatter *PlainTextFormatter
	// 默认时间戳格式
	formatter = &PlainTextFormatter{}
	formatter.initialize()
	if formatter.TimestampFormat != time.RFC3339 {
		t.Error("unexpected default timestamp format:", formatter.TimestampFormat)
	}
	// 测试事件字段
	formatter = &PlainTextFormatter{
		EventFormat: "%(.event_fields_space_seperated_text|s)",
	}
	formatter.initialize()
	if !formatter.needEventFieldsSpaceSeperatedText {
		t.Error("needEventFieldsSpaceSeperatedText is false")
	}
	formatter = &PlainTextFormatter{
		EventFormat: "%(.event_fields_comma_seperated_text|s)",
	}
	formatter.initialize()
	if !formatter.needEventFieldsCommaSeperatedText {
		t.Error("needEventFieldsCommaSeperatedText is false")
	}
	// 测试会话字段
	formatter = &PlainTextFormatter{
		EventFormat: "%(.session_fields_space_seperated_text|s)",
	}
	formatter.initialize()
	if !formatter.needSessionFieldsSpaceSeperatedText {
		t.Error("needSessionFieldsSpaceSeperatedText is false")
	}
	formatter = &PlainTextFormatter{
		EventFormat: "%(.session_fields_comma_seperated_text|s)",
	}
	formatter.initialize()
	if !formatter.needSessionFieldsCommaSeperatedText {
		t.Error("needSessionFieldsCommaSeperatedText is false")
	}
	// 测试全局字段
	formatter = &PlainTextFormatter{
		EventFormat: "%(.global_fields_space_seperated_text|s)",
	}
	formatter.initialize()
	if !formatter.needGlobalFieldsSpaceSeperatedText {
		t.Error("needGlobalFieldsSpaceSeperatedText is false")
	}
	formatter = &PlainTextFormatter{
		EventFormat: "%(.global_fields_comma_seperated_text|s)",
	}
	formatter.initialize()
	if !formatter.needGlobalFieldsCommaSeperatedText {
		t.Error("needGlobalFieldsCommaSeperatedText is false")
	}
	// 测试全局字段
	formatter = &PlainTextFormatter{
		EventFormat: "%(.all_fields_space_seperated_text|s)",
	}
	formatter.initialize()
	if !formatter.needAllFieldsSpaceSeperatedText {
		t.Error("needAllFieldsSpaceSeperatedText is false")
	}
	formatter = &PlainTextFormatter{
		EventFormat: "%(.all_fields_comma_seperated_text|s)",
	}
	formatter.initialize()
	if !formatter.needAllFieldsCommaSeperatedText {
		t.Error("needAllFieldsCommaSeperatedText is false")
	}
	// 测试Level颜色构造
	formatter = &PlainTextFormatter{
		LevelANSIColors: map[string]string{
			"info":  "FgHiBlue",
			"warn":  "FgHiYellow",
			"error": "FgHiRed",
		},
	}
	formatter.initialize()
	if formatter.levelANSIColorFuncs["info"] == nil ||
		formatter.levelANSIColorFuncs["warn"] == nil ||
		formatter.levelANSIColorFuncs["error"] == nil {
		t.Error("unexpected level ansi color funcs:", formatter.levelANSIColorFuncs)
	}
}

func TestPlainTextFormatterFieldify(t *testing.T) {
	globalFields = map[string]interface{}{
		"zero": 0,
		"one":  1,
	}
	session := NewSession().WithFields(Fields{"foo": "bar", "bar": "foo"})
	event := session.Event().WithFields(Fields{"true": true, "false": false})
	formatter := &PlainTextFormatter{
		EventFormat: strings.Join([]string{
			"%(.event_fields_space_seperated_text|s)",
			"%(.event_fields_comma_seperated_text|s)",
			"%(.session_fields_space_seperated_text|s)",
			"%(.session_fields_comma_seperated_text|s)",
			"%(.global_fields_space_seperated_text|s)",
			"%(.global_fields_comma_seperated_text|s)",
			"%(.all_fields_space_seperated_text|s)",
			"%(.all_fields_comma_seperated_text|s)",
		}, "\n"),
		LevelANSIColors: map[string]string{
			"info": "FgHiBlue",
		},
		SortFields: true,
	}
	formatter.initialize()
	event.Level = "info"
	fields := formatter.fieldify(event)
	if text, ok := fields[".event_fields_space_seperated_text"].(string); !ok || text != "false=false true=true" {
		t.Error("unexpected event_fields_space_seperated_text value:", fields[".event_fields_space_seperated_text"])
	}
	if text, ok := fields[".event_fields_comma_seperated_text"].(string); !ok || text != "false=false,true=true" {
		t.Error("unexpected event_fields_comma_seperated_text value:", fields[".event_fields_comma_seperated_text"])
	}
	if text, ok := fields[".session_fields_space_seperated_text"].(string); !ok || text != "bar=\"foo\" foo=\"bar\"" {
		t.Error("unexpected session_fields_space_seperated_text value:", fields[".session_fields_space_seperated_text"])
	}
	if text, ok := fields[".session_fields_comma_seperated_text"].(string); !ok || text != "bar=\"foo\",foo=\"bar\"" {
		t.Error("unexpected session_fields_comma_seperated_text value:", fields[".session_fields_comma_seperated_text"])
	}
	if text, ok := fields[".global_fields_space_seperated_text"].(string); !ok || text != "one=1 zero=0" {
		t.Error("unexpected global_fields_space_seperated_text value:", fields[".global_fields_space_seperated_text"])
	}
	if text, ok := fields[".global_fields_comma_seperated_text"].(string); !ok || text != "one=1,zero=0" {
		t.Error("unexpected global_fields_comma_seperated_text value:", fields[".global_fields_comma_seperated_text"])
	}
	if text, ok := fields[".all_fields_space_seperated_text"].(string); !ok || text != "one=1 zero=0 bar=\"foo\" foo=\"bar\" false=false true=true" {
		t.Error("unexpected all_fields_space_seperated_text value:", fields[".all_fields_space_seperated_text"])
	}
	if text, ok := fields[".all_fields_comma_seperated_text"].(string); !ok || text != "one=1,zero=0,bar=\"foo\",foo=\"bar\",false=false,true=true" {
		t.Error("unexpected all_fields_comma_seperated_text value:", fields[".all_fields_comma_seperated_text"])
	}
	if fields["caller.file"] != "formatter_test.go" || fields["caller.func"] != "TestPlainTextFormatterFieldify" || fields["caller.package"] != "github.com/yangchenxing/go-slog" {
		t.Errorf("unexpected caller: package=%v, file=%v, func=%v\n", fields["caller.package"], fields["caller.file"], fields["caller.func"])
	}
	levelColorFunc := color.New(color.FgHiBlue).SprintFunc()
	if fields["level"] != levelColorFunc("info") {
		t.Error("unexpected level:", fields["level"])
	}
}

func TestPlainTextFormatterFormatEvent(t *testing.T) {
	formatter := &PlainTextFormatter{
		EventFormat: "%(caller.package|s)/%(caller.file|s):%(caller.func|s)",
	}
	event := newEvent(1, nil)
	content, err := formatter.FormatEvent(event)
	if err != nil {
		t.Error("format fail:", err.Error())
	} else if string(content) != "github.com/yangchenxing/go-slog/formatter_test.go:TestPlainTextFormatterFormatEvent" {
		t.Error("unexpected result:", string(content))
	}
}

func TestPlainTextFormatterFormatEvents(t *testing.T) {
	formatter := &PlainTextFormatter{
		EventFormat: "%(caller.package|s)/%(caller.file|s):%(caller.func|s) %(sample|d)",
	}
	content, err := formatter.FormatEvents([]*Event{
		newEvent(1, nil).WithField("sample", 1),
		newEvent(1, nil).WithField("sample", 2)})
	expected := strings.Join([]string{
		"github.com/yangchenxing/go-slog/formatter_test.go:TestPlainTextFormatterFormatEvents 1",
		"github.com/yangchenxing/go-slog/formatter_test.go:TestPlainTextFormatterFormatEvents 2",
	}, "\n") + "\n"
	if err != nil {
		t.Error("format fail:", err.Error())
	} else if string(content) != expected {
		t.Errorf("unexpected result: actual=%q, expected=%q\n", string(content), expected)
	}
}
