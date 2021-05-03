package errors

import (
	"fmt"

	"github.com/payfazz/go-errors/v2/trace"
)

func new(skip int, data interface{}, err error, deep int) error {
	ret := &Error{
		Data:       data,
		Cause:      err,
		StackTrace: trace.Get(skip+1, deep),
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
func New(data interface{}) error {
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

// ErrorWithCauseF is same with like ErrorF, but you can spesify the error cause
func ErrorWithCauseF(cause error, f string, v ...interface{}) error {
	return new(1, fmt.Sprintf(f, v...), cause, DefaultDeep)
}

// NewWithCauseAndDeep is same with NewWithCause, but with specified stack deep
func NewWithCauseAndDeep(data interface{}, err error, deep int) error {
	return new(1, data, err, deep)
}
