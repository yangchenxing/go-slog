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

	WriterFactory.RegisterInstance("stdout", StdoutWriter)
	WriterFactory.RegisterInstance("stderr", StderrWriter)
	WriterFactory.RegisterType("email", reflect.TypeOf((*PlainTextEmailHandler)(nil)).Elem())
	WriterFactory.RegisterType("time_rotated_file", reflect.TypeOf((*TimeRotatedFileWriter)(nil)).Elem())

	map2struct.RegisterFactory(HandlerFactory)
	map2struct.RegisterFactory(WriterFactory)
}
