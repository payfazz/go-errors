package errors

import (
	"fmt"

	"github.com/payfazz/go-errors/v2/trace"
)

func New(text string) error {
	return newErrorf(1, "%s", text)
}

func Errorf(format string, a ...interface{}) error {
	return newErrorf(1, format, a...)
}

type errorf struct {
	e error
	l []trace.Location
}

func newErrorf(skip int, format string, a ...interface{}) error {
	return &errorf{
		e: fmt.Errorf(format, a...),
		l: trace.Get(skip+1, DefaultTraceDeep),
	}
}

func (t *errorf) Error() string {
	return t.e.Error()
}

func (t *errorf) Unwrap() error {
	return Unwrap(t.e)
}

func (t *errorf) StackTrace() []trace.Location {
	return t.l
}
