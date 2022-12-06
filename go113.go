//go:build !go1.20

package errors

import "github.com/payfazz/go-errors/v2/trace"

func newTraced(err error) error {
	return &tracedErr{err, trace.Get(2, traceDeep)}
}
