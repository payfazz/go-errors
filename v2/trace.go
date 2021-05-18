package errors

import (
	"github.com/payfazz/go-errors/v2/trace"
)

type tracedErr struct {
	err    error
	trace  []trace.Location
	parent []trace.Location
}

func (e *tracedErr) Error() string {
	return e.err.Error()
}

func (e *tracedErr) Unwrap() error {
	return Unwrap(e.err)
}

func (e *tracedErr) As(target interface{}) bool {
	return As(e.err, target)
}

func (e *tracedErr) Is(target error) bool {
	return Is(e.err, target)
}

const defaultDeep = 150

// Trace will return new error that have stack trace
//
// will return same err if err already have stack trace
func Trace(err error) error {
	if err == nil {
		return nil
	}

	if _, ok := err.(*tracedErr); ok {
		return err
	}

	return &tracedErr{
		err:   err,
		trace: trace.Get(1, defaultDeep),
	}
}

// Get stack trace of err
//
// return nil if err doesn't have stack trace
func StackTrace(err error) []trace.Location {
	if e, ok := err.(*tracedErr); ok {
		return e.trace
	}
	return nil
}
