package errors

import (
	"fmt"
	"strings"
)

// Error represent the wrapped error
type Error struct {
	text       string
	cause      error
	stackTrace []Location
}

var _ error = (*Error)(nil)

func (e *Error) Error() string {
	if e.text != "" {
		return e.text
	}
	if e.cause != nil {
		return e.cause.Error()
	}
	return ""
}

// Cause return the error that cause this error
func (e *Error) Cause() error {
	return e.cause
}

// StackTrace retrun the stack trace when this error created
func (e *Error) StackTrace() []Location {
	return e.stackTrace
}

// String representation of Error
func (e *Error) String() string {
	var buff strings.Builder
	var cause error = e
	var first = true

	for cause != nil {
		if first {
			first = false
		} else {
			buff.WriteString("Caused by ")
		}
		buff.WriteString("Error: ")
		buff.WriteString(cause.Error())
		buff.WriteString("\n")
		if err, ok := cause.(*Error); ok {
			for _, l := range err.StackTrace() {
				buff.WriteString("- ")
				buff.WriteString(l.String())
				buff.WriteString("\n")
			}
			cause = err.Cause()
		} else {
			cause = nil
		}
	}

	return buff.String()
}

func new(skip int, text string, err error, deep int) error {
	ret := &Error{
		text:       text,
		cause:      err,
		stackTrace: generateStackTrace(skip+1, deep),
	}

	return ret
}

func wrap(skip int, text string, err error, deep int) error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*Error); ok {
		return e
	}
	return new(skip+1, text, err, deep)
}

// Wrap the err, if err is nil, then return nil
func Wrap(err error) error {
	return wrap(1, "", err, DefaultDeep)
}

// WrapWithDeep is same with Wrap, but with specified stack deep
func WrapWithDeep(err error, deep int) error {
	return wrap(1, "", err, deep)
}

// New return an Error with the given text.
func New(text string) error {
	return new(1, text, nil, DefaultDeep)
}

// Errorf return an Error with text according to a format specifier.
//
// Format as specified by fmt.Sprintf
func Errorf(f string, v ...interface{}) error {
	return new(1, fmt.Sprintf(f, v...), nil, DefaultDeep)
}

// ErrorfWithDeep return an Error with text according to a format specifier.
//
// Format as specified by fmt.Sprintf
func ErrorfWithDeep(deep int, f string, v ...interface{}) error {
	return new(1, fmt.Sprintf(f, v...), nil, deep)
}

// NewWithDeep is same with New, but with specified stack deep
func NewWithDeep(text string, deep int) error {
	return new(1, text, nil, deep)
}

// NewWithCause is same with New, but it also indicate that this Error is caused by err.
func NewWithCause(text string, err error) error {
	return new(1, text, err, DefaultDeep)
}

// NewWithCauseAndDeep is same with NewWithCause, but with specified stack deep
func NewWithCauseAndDeep(text string, err error, deep int) error {
	return new(1, text, err, deep)
}

// Format the error as string
func Format(err error) string {
	if err == nil {
		return ""
	}

	if err2, ok := err.(*Error); ok {
		return err2.String()
	}

	return err.Error()
}

// Cause return the root cause of the error
func Cause(err error) error {
	last := err
	for {
		if err == nil {
			return last
		}

		if err2, ok := err.(*Error); ok {
			last = err
			err = err2.Cause()
		} else {
			return err
		}
	}
}
