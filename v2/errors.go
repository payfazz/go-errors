// Package errors.
//
// This package can be used as drop-in replacement for https://golang.org/pkg/errors package.
//
// This package provide StackTrace function to get the stack trace,
// use Wrap to make sure the error have stack trace
package errors

import (
	stderrors "errors"
	"fmt"

	"github.com/payfazz/go-errors/v2/trace"
)

const defaultDeep = 20

type StackTracer interface {
	StackTrace() []trace.Location
}

// see https://golang.org/pkg/errors/#As
func As(err error, target interface{}) bool {
	return stderrors.As(err, target)
}

// see https://golang.org/pkg/errors/#Is
func Is(err, target error) bool {
	return stderrors.Is(err, target)
}

// see https://golang.org/pkg/errors/#Unwrap
func Unwrap(err error) error {
	return stderrors.Unwrap(err)
}

// Ignore the err, useful for static code analysis
func Ignore(err error) {}

// see https://golang.org/pkg/errors/#New
func New(text string) error {
	return newTracedErr(1, defaultDeep, &textErr{text: text})
}

// like New, but you can specify the stack trace deep
func NewWithDeep(deep int, text string) error {
	return newTracedErr(1, deep, &textErr{text: text})
}

// like New, but you can specify the cause error
func NewWithCause(text string, cause error) error {
	return newTracedErr(1, defaultDeep, &textErr{text: text, cause: cause})
}

// like NewWithCause, but you can specify the stack trace deep
func NewWithCauseAndDeep(deep int, text string, cause error) error {
	return newTracedErr(1, deep, &textErr{text: text, cause: cause})
}

// see https://golang.org/pkg/fmt/#Errorf
func Errorf(format string, a ...interface{}) error {
	err := fmt.Errorf(format, a...)
	cause := Unwrap(err)
	return newTracedErr(1, defaultDeep, &textErr{text: err.Error(), cause: cause})
}

// like Errorf, but you can specify the stack trace deep
func ErrorfWithDeep(deep int, format string, a ...interface{}) error {
	err := fmt.Errorf(format, a...)
	cause := Unwrap(err)
	return newTracedErr(1, deep, &textErr{text: err.Error(), cause: cause})
}

// Wrap the err if the err doens't have stack trace
func Wrap(err error) error {
	return newTracedErr(1, defaultDeep, err)
}

// like Wrap, but you can specify the stack trace deep
func WrapWithDeep(deep int, err error) error {
	return newTracedErr(1, deep, err)
}

// Get stack trace of where the error is generated, return nil if none
func StackTrace(err error) []trace.Location {
	if e, ok := err.(StackTracer); ok {
		return e.StackTrace()
	}
	return nil
}
