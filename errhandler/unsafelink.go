package errhandler

import _ "unsafe"

//go:linkname errors_wrap github.com/payfazz/go-errors.wrap
func errors_wrap(skip int, text string, err error, deep int) error
