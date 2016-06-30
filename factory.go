package slog

import (
	"reflect"

	"github.com/yangchenxing/go-map2struct"
)

var (
	HandlerFactory = map2struct.NewGeneralInterfaceFactory(reflect.TypeOf((*Handler)(nil)).Elem(), "type", nil)
	WriterFactory  = map2struct.NewGeneralInterfaceFactory(reflect.TypeOf((*Writer)(nil)).Elem(), "type", nil)
)

func init() {
	HandlerFactory.RegisterType("json", reflect.TypeOf((*JsonHandler)(nil)).Elem())
	HandlerFactory.RegisterType("plaintext", reflect.TypeOf((*PlainTextHandler)(nil)).Elem())

	HandlerFactory.RegisterInstance("stdout", StdoutWriter)
	HandlerFactory.RegisterInstance("stderr", StderrWriter)
	HandlerFactory.RegisterType("email", reflect.TypeOf((*PlainTextEmailHandler)(nil)).Elem())
	HandlerFactory.RegisterType("time_rotated_file", reflect.TypeOf((*TimeRotatedFileWriter)(nil)).Elem())
}
