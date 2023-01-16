//go:build go1.20

package errors

import (
	stderrors "errors"

	"github.com/payfazz/go-errors/v2/trace"
)

func newTraced(err error) error {
	if _, ok := err.(unwrapslice); ok {
		return &tracedErrSlice{tracedErr{err, trace.Get(2, traceDeep)}}
	}
	return &tracedErr{err, trace.Get(2, traceDeep)}
}

type unwrapslice interface{ Unwrap() []error }

type tracedErrSlice struct{ tracedErr }

func (e *tracedErrSlice) Unwrap() []error {
	return e.err.(unwrapslice).Unwrap()
}

// see [stdlib errors.Join]
//
// [stdlib errors.Join]: https://pkg.go.dev/errors/#As
func Join(err ...error) error {
	return newTraced(stderrors.Join(err...))
}
