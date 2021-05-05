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

const defaultDeep = 50

// Wrap the err if the err doens't have stack trace
func Wrap(err error) error {
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

// Get stack trace of where the error is generated or wrapped, return nil if none
func StackTrace(err error) []trace.Location {
	if e, ok := err.(*tracedErr); ok {
		return e.trace
	}
	return nil
}
