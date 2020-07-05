package errors

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/payfazz/go-errors/trace"
)

// Error represent the wrapped error
//
// all error returned from New*, Errorf, Wrap will have type *Error
type Error struct {
	// Data is arbitrary data attached to this error
	Data interface{}

	// Cause of this error
	Cause error

	// StackTrace where this error generated
	StackTrace []trace.Location
}

// Error from error interface
//
// will return the error message
func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if data := e.Data; data != nil {
		return fmt.Sprint(data)
	}
	if cause := e.Cause; cause != nil {
		return cause.Error()
	}
	return ""
}

// Walk is helper for walk the error chains.
//
// Walk will return true if fn is called at least once
//
// Walk will walk the chains as long fn return true
func Walk(err error, fn func(*Error) bool) bool {
	executedOnce := false
	for err != nil {
		wrappedErr, isWrapped := err.(*Error)
		if isWrapped {
			executedOnce = true
			if !fn(wrappedErr) {
				break
			}
			err = wrappedErr.Cause // next error in the chain
		} else {
			break
		}
	}
	return executedOnce
}

// Format representation of the Error, including stack trace.
//
// Use err.Error() if you want to get just the error string
func Format(err error) string {
	return format(err, -1)
}

// FormatWithDeep representation of the Error, including stack trace with specified deep.
//
// Use err.Error() if you want to get just the error string
func FormatWithDeep(err error, deep int) string {
	return format(err, deep)
}

// RootCause return the root cause of the error
func RootCause(err error) error {
	Walk(err, func(wrappedErr *Error) bool {
		if cause := wrappedErr.Cause; cause != nil {
			err = cause
		} else {
			err = wrappedErr
		}
		return true
	})
	return err
}

// InErrorChain will return true if the data is exists in error chain
func InErrorChain(err error, data interface{}) bool {
	if !reflect.TypeOf(data).Comparable() {
		return false
	}

	if err == data {
		return true
	}

	match := false
	Walk(err, func(err *Error) bool {
		if err == data || err.Data == data || err.Cause == data {
			match = true
			return false
		}
		return true
	})
	return match
}

// Ignore the error.
func Ignore(err error) {}

func format(err error, deep int) string {
	if err == nil {
		return ""
	}

	var buff strings.Builder
	var first = true
	if Walk(err, func(wrappedErr *Error) bool {
		printErr := func(err error) {
			if first {
				first = false
			} else {
				buff.WriteString("Caused by ")
			}
			buff.WriteString("Error: ")
			buff.WriteString(err.Error())
		}
		printErr(wrappedErr)
		buff.WriteString("\n")
		for i, l := range wrappedErr.StackTrace {
			if i == deep {
				break
			}
			buff.WriteString("- ")
			buff.WriteString(l.String())
			buff.WriteString("\n")
		}

		// Walk will skip this last one, so we need to do it here
		if wrappedErr.Cause != nil {
			if _, ok := wrappedErr.Cause.(*Error); !ok {
				printErr(wrappedErr.Cause)
			}
		}

		return true
	}) {
		return buff.String()
	}

	return "Error: " + err.Error()
}
