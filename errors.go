package errors

import (
	stderrors "errors"
	"fmt"

	"github.com/payfazz/go-errors/v2/trace"
)

// see https://pkg.go.dev/errors/#As
func As(err error, target interface{}) bool {
	return stderrors.As(err, target)
}

// see https://pkg.go.dev/errors/#Is
func Is(err, target error) bool {
	return stderrors.Is(err, target)
}

// see https://pkg.go.dev/errors/#Unwrap
func Unwrap(err error) error {
	return stderrors.Unwrap(err)
}

// see https://pkg.go.dev/errors/#New
func New(text string) error {
	return &traced{stderrors.New(text), trace.Get(1, traceDeep)}
}

// see https://pkg.go.dev/fmt/#Errorf
func Errorf(format string, a ...interface{}) error {
	return &traced{fmt.Errorf(format, a...), trace.Get(1, traceDeep)}
}

type wrapper struct {
	msg   string
	cause error
}

func (w *wrapper) Error() string { return w.msg }
func (w *wrapper) Unwrap() error { return w.cause }

// like New, but you can specify the cause error
func NewWithCause(text string, cause error) error {
	return &traced{&wrapper{text, cause}, trace.Get(1, traceDeep)}
}
