package errors

import (
	"strings"
)

// Error represent the wrapped error
type Error interface {
	// Cause return the error that cause this error
	Cause() error

	// StackTrace retrun the stack trace when this error created
	StackTrace() []Location

	// String representation of Error
	String() string

	error

	// internal is just empty function, the purpose is to make this interface cannot be implemented outside this package
	internal()
}

type errorType struct {
	text       string
	cause      error
	stackTrace []Location
}

func (e *errorType) Error() string {
	if e.text != "" {
		return e.text
	}
	if e.cause != nil {
		return e.cause.Error()
	}
	return ""
}

func (e *errorType) Cause() error {
	return e.cause
}

func (e *errorType) StackTrace() []Location {
	return e.stackTrace
}

func (e *errorType) String() string {
	var buff strings.Builder
	var cause error = e
	var first = true

	for cause != nil {
		if first {
			first = false
		} else {
			buff.WriteString("\nCaused by ")
		}
		buff.WriteString("Error: ")
		buff.WriteString(cause.Error())
		buff.WriteString("\n")
		if err, ok := cause.(Error); ok {
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

func (e *errorType) internal() {}

func new(skip int, text string, err error, deep int) Error {
	ret := &errorType{
		text:       text,
		cause:      err,
		stackTrace: generateStackTrace(skip+1, deep),
	}

	return ret
}

func wrap(skip int, text string, err error, deep int) Error {
	if err == nil {
		return nil
	}
	if e, ok := err.(Error); ok {
		return e
	}
	return new(skip+1, text, err, deep)
}

// Wrap the err, if err is nil, then return nil
func Wrap(err error) Error {
	return wrap(1, "", err, defaultDeep)
}

// WrapWithDeep is same with Wrap, but with specified stack deep
func WrapWithDeep(err error, deep int) Error {
	return wrap(1, "", err, deep)
}

// New returns an Error that formats as the given text.
func New(text string) Error {
	return new(1, text, nil, defaultDeep)
}

// NewWithDeep is same with New, but with specified stack deep
func NewWithDeep(text string, deep int) Error {
	return new(1, text, nil, deep)
}

// NewWithCause returns an Error that formats as the given text,
// it also indicate that this Error is caused by err.
func NewWithCause(text string, err error) Error {
	return new(1, text, err, defaultDeep)
}

// NewWithCauseAndDeep is same with NewWithCause, but with specified stack deep
func NewWithCauseAndDeep(text string, err error, deep int) Error {
	return new(1, text, err, deep)
}

// Format the error as string
func Format(err error) string {
	if err == nil {
		return ""
	}

	if err2, ok := err.(Error); ok {
		return err2.String()
	}

	return err.Error()
}

// RealCause is same with Cause.
//
// Deprecated: use Cause
func RealCause(err error) error {
	return Cause(err)
}

// Cause return the root cause of the error
func Cause(err error) error {
	last := err
	for {
		if err == nil {
			return last
		}

		if err2, ok := err.(Error); ok {
			last = err
			err = err2.Cause()
		} else {
			return err
		}
	}
}
