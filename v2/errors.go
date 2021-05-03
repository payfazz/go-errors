package errors

import (
	stderrors "errors"

	"github.com/payfazz/go-errors/v2/trace"
)

const DefaultTraceDeep = 20

func As(err error, target interface{}) bool {
	return stderrors.As(err, target)
}

func Is(err, target error) bool {
	return stderrors.Is(err, target)
}

func Unwrap(err error) error {
	return stderrors.Unwrap(err)
}

func StackTrace(err error) []trace.Location {
	if e, ok := err.(interface{ StackTrace() []trace.Location }); ok {
		return e.StackTrace()
	}

	return nil
}
