package errors

import (
	"github.com/payfazz/go-errors/v2/trace"
)

type tracedErr struct {
	error
	trace []trace.Location
}

func newTracedErr(skip int, deep int, err error) error {
	if err == nil {
		return nil
	}

	if _, ok := err.(StackTracer); ok {
		return err
	}

	return &tracedErr{
		error: err,
		trace: trace.Get(skip+1, deep),
	}
}

func (e *tracedErr) Unwrap() error {
	return Unwrap(e.error)
}

func (e *tracedErr) As(target interface{}) bool {
	return As(e.error, target)
}

func (e *tracedErr) Is(target error) bool {
	return Is(e.error, target)
}

func (e *tracedErr) StackTrace() []trace.Location {
	return e.trace
}
