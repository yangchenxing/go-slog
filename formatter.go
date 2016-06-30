package slog

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"sort"

	"github.com/yangchenxing/go-string-mapformatter"
)

type EventFormatter interface {
	FormatEvent(*Event) ([]byte, error)
}

type MultiEventFormatter interface {
	FormatEvents([]*Event) ([]byte, error)
}

type PlainTextFormatter struct {
	TimestampFormat                     string
	EventFormat                         string
	SortFields                          bool
	levelANSIColorFuncs                 map[string]func(...interface{}) string
	needEventFieldsSpaceSeperatedText   bool
	needEventFieldsCommaSeperatedText   bool
	needSessionFieldsSpaceSeperatedText bool
	needSessionFieldsCommaSeperatedText bool
	needGlobalFieldsSpaceSeperatedText  bool
	needGlobalFieldsCommaSeperatedText  bool
	needAllFieldsSpaceSeperatedText     bool
	needAllFieldsCommaSeperatedText     bool
	initialized                         bool
}

func (formatter *PlainTextFormatter) FormatEvents(events []*Event) ([]byte, error) {
	if !formatter.initialized {
		formatter.initialize()
	}
	var buffer bytes.Buffer
	for _, event := range events {
		fields := formatter.fieldify(event)
		buffer.Write([]byte(mapformatter.Format(formatter.EventFormat, fields)))
		buffer.WriteRune('\n')
	}
	return buffer.Bytes(), nil
}

func (formatter *PlainTextFormatter) FormatEvent(event *Event) ([]byte, error) {
	if !formatter.initialized {
		formatter.initialize()
	}
	fields := formatter.fieldify(event)
	content := mapformatter.Format(formatter.EventFormat, fields)
	return []byte(content), nil
}

func (formatter *PlainTextFormatter) fieldify(event *Event) map[string]interface{} {
	fields := event.Fieldify(formatter.TimestampFormat)
	// 补充事件自定义字段连接文本
	if formatter.needEventFieldsSpaceSeperatedText {
		fields[".event_fields_space_seperated_text"] = JoinFields(event.Fields, "=", " ", formatter.SortFields)
	}
	if formatter.needEventFieldsCommaSeperatedText {
		fields[".event_fields_comma_seperated_text"] = JoinFields(event.Fields, "=", ",", formatter.SortFields)
	}
	// 补充会话自定义字段连接文本
	if formatter.needSessionFieldsSpaceSeperatedText {
		fields[".session_fields_space_seperated_text"] = JoinFields(event.Session, "=", " ", formatter.SortFields)
	}
	if formatter.needSessionFieldsCommaSeperatedText {
		fields[".session_fields_comma_seperated_text"] = JoinFields(event.Session, "=", ",", formatter.SortFields)
	}
	// 补充全局自定义字段连接文本
	if formatter.needSessionFieldsSpaceSeperatedText {
		fields[".global_fields_space_seperated_text"] = JoinFields(globalFields, "=", " ", formatter.SortFields)
	}
	if formatter.needSessionFieldsCommaSeperatedText {
		fields[".global_fields_comma_seperated_text"] = JoinFields(globalFields, "=", ",", formatter.SortFields)
	}
	// 补充全字段连接文本
	if formatter.needAllFieldsSpaceSeperatedText {
		parts := make([]string, 0, 3)
		if text := JoinFields(globalFields, "=", " ", formatter.SortFields); text != "" {
			parts = append(parts, text)
		}
		if text := JoinFields(event.Session, "=", " ", formatter.SortFields); text != "" {
			parts = append(parts, text)
		}
		if text := JoinFields(event.Fields, "=", " ", formatter.SortFields); text != "" {
			parts = append(parts, text)
		}
		fields[".all_fields_space_seperated_text"] = strings.Join(parts, " ")
	}
	if formatter.needAllFieldsCommaSeperatedText {
		parts := make([]string, 0, 3)
		if text := JoinFields(globalFields, "=", ",", formatter.SortFields); text != "" {
			parts = append(parts, text)
		}
		if text := JoinFields(event.Session, "=", ",", formatter.SortFields); text != "" {
			parts = append(parts, text)
		}
		if text := JoinFields(event.Fields, "=", ",", formatter.SortFields); text != "" {
			parts = append(parts, text)
		}
		fields[".all_fields_comma_seperated_text"] = strings.Join(parts, ",")
	}
	// Level字段着色
	if colorFunc := formatter.levelANSIColorFuncs[event.Level]; colorFunc != nil {
		fields["level"] = colorFunc(event.Level)
	} else {
		fields["level"] = event.Level
	}
	// Caller字段打平
	fields["caller.file"] = event.Caller.File
	fields["caller.func"] = event.Caller.Func
	fields["caller.line"] = event.Caller.Line
	fields["caller.package"] = event.Caller.Package
	return fields
}

func JoinFields(fields map[string]interface{}, equal, seperator string, order bool) string {
	strs := make([]string, len(fields))
	i := 0
	for key, value := range fields {
		if s, ok := value.(string); ok {
			strs[i] = fmt.Sprintf("%s%s%q", key, equal, s)
		} else {
			strs[i] = fmt.Sprintf("%s%s%v", key, equal, value)
		}
		i++
	}
	if order {
		sort.Sort(sort.StringSlice(strs))
	}
	return strings.Join(strs, seperator)
}

func (formatter *PlainTextFormatter) initialize() {
	// 补充默认时间戳格式
	if formatter.TimestampFormat == "" {
		formatter.TimestampFormat = time.RFC3339
	}
	// 判断是否需要使用自定义Field字段，非准确判断
	formatter.needEventFieldsSpaceSeperatedText = strings.Contains(formatter.EventFormat, "%(.event_fields_space_seperated_text|")
	formatter.needEventFieldsCommaSeperatedText = strings.Contains(formatter.EventFormat, "%(.event_fields_comma_seperated_text|")
	formatter.needSessionFieldsSpaceSeperatedText = strings.Contains(formatter.EventFormat, "%(.session_fields_space_seperated_text|")
	formatter.needSessionFieldsCommaSeperatedText = strings.Contains(formatter.EventFormat, "%(.session_fields_comma_seperated_text|")
	formatter.needGlobalFieldsSpaceSeperatedText = strings.Contains(formatter.EventFormat, "%(.global_fields_space_seperated_text|")
	formatter.needGlobalFieldsCommaSeperatedText = strings.Contains(formatter.EventFormat, "%(.global_fields_comma_seperated_text|")
	formatter.needAllFieldsSpaceSeperatedText = strings.Contains(formatter.EventFormat, "%(.all_fields_space_seperated_text|")
	formatter.needAllFieldsCommaSeperatedText = strings.Contains(formatter.EventFormat, "%(.all_fields_comma_seperated_text|")
	formatter.initialized = true
}
