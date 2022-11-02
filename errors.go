package errors

import (
	stderrors "errors"
	"fmt"

	"github.com/payfazz/go-errors/v2/trace"
)

// see [errors.As]
func As(err error, target any) bool {
	return stderrors.As(err, target)
}

// see [errors.Is]
func Is(err, target error) bool {
	return stderrors.Is(err, target)
}

// see [errors.Unwrap]
func Unwrap(err error) error {
	return stderrors.Unwrap(err)
}

// see [errors.New]
func New(text string) error {
	return &traced{stderrors.New(text), trace.Get(1, traceDeep)}
}

// see [fmt.Errorf]
func Errorf(format string, a ...any) error {
	return &traced{fmt.Errorf(format, a...), trace.Get(1, traceDeep)}
}

// will panic if err is not nil
func Check(err error) {
	if err != nil {
		panic(Trace(err))
	}
}

// Assert will panic if fact is false
func Assert(fact bool) {
	if !fact {
		panic(New("assertion failed"))
	}
}
