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
		err:   &anyErr{data: text},
		trace: trace.Get(1, defaultDeep),
	}
}

// like New, but you can specify the cause error
func NewWithCause(text string, cause error) error {
	return &tracedErr{
		err:   &anyErr{data: text, cause: cause},
		trace: trace.Get(1, defaultDeep),
	}
}
