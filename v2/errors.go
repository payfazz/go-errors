// Package errors.
//
// This package can be used as drop-in replacement for https://golang.org/pkg/errors package.
//
// This package provide StackTrace function to get the stack trace,
// use Wrap to make sure the error have stack trace
package errors

import (
	stderrors "errors"

	"github.com/payfazz/go-errors/v2/trace"
)

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

// see https://golang.org/pkg/errors/#New
func New(text string) error {
	return &tracedErr{
		error: &anyErr{data: text},
		trace: trace.Get(1, defaultDeep),
	}
}

// like New, but you can specify the cause error
func NewWithCause(text string, cause error) error {
	return &tracedErr{
		error: &anyErr{data: text, cause: cause},
		trace: trace.Get(1, defaultDeep),
	}
}
