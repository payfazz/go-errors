package errors

import (
	"fmt"
	"strings"

	"github.com/payfazz/go-errors/trace"
)

// Error represent the wrapped error
type Error struct {
	data       interface{}
	cause      error
	stackTrace []trace.Location
}

var _ error = (*Error)(nil)

func (e *Error) Error() string {
	if e.data != nil {
		switch v := e.data.(type) {
		case string:
			return v
		case error:
			return v.Error()
		default:
			return fmt.Sprint(v)
		}
	}
	if e.cause != nil {
		return e.cause.Error()
	}
	return ""
}

// Data return the error data
func (e *Error) Data() interface{} {
	return e.data
}

// Cause return the error that cause this error
func (e *Error) Cause() error {
	return e.cause
}

// StackTrace retrun the stack trace when this error created
func (e *Error) StackTrace() []trace.Location {
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

func new(skip int, data interface{}, err error, deep int) error {
	ret := &Error{
		data:       data,
		cause:      err,
		stackTrace: trace.Get(skip+1, deep),
	}

	return ret
}

func wrap(skip int, data interface{}, err error, deep int) error {
	if e, ok := err.(*Error); ok {
		return e
	}
	return new(skip+1, data, err, deep)
}

// Wrap the err, if err is nil, then return nil
func Wrap(err error) error {
	if err == nil {
		return nil
	}
	return wrap(1, nil, err, DefaultDeep)
}

// WrapWithDeep is same with Wrap, but with specified stack deep
func WrapWithDeep(err error, deep int) error {
	if err == nil {
		return nil
	}
	return wrap(1, nil, err, deep)
}

// New return an Error with the given data.
func New(data string) error {
	return new(1, data, nil, DefaultDeep)
}

// Errorf return an Error with text according to a format specifier.
//
// Format as specified by fmt.Sprintf
func Errorf(f string, v ...interface{}) error {
	return new(1, fmt.Sprintf(f, v...), nil, DefaultDeep)
}

// NewWithCause is same with New, but it also indicate that this Error is caused by err.
func NewWithCause(data interface{}, err error) error {
	return new(1, data, err, DefaultDeep)
}

// NewWithCauseAndDeep is same with NewWithCause, but with specified stack deep
func NewWithCauseAndDeep(data interface{}, err error, deep int) error {
	return new(1, data, err, deep)
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

// InErrorChain check if data in error chain
func InErrorChain(err error, data interface{}) bool {
	for {
		if err == nil {
			return false
		}

		if err == data {
			return true
		}

		if err2, ok := err.(*Error); ok {
			if err2.Data() == data {
				return true
			}
			err = err2.Cause()
		} else {
			return false
		}
	}
}
