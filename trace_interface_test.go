package errors_test

import (
	"testing"

	"github.com/payfazz/go-errors/v2"
	"github.com/payfazz/go-errors/v2/trace"
)

type myTracedError struct{ trace []trace.Location }

func (m *myTracedError) Error() string                { return "" }
func (m *myTracedError) StackTrace() []trace.Location { return m.trace }

func TestStackTraceInterface(t *testing.T) {
	var err error
	funcAA(func() {
		funcBB(func() {
			err = &myTracedError{trace.Get(0, 100)}
		})
	})

	trace := errors.StackTrace(err)

	if !haveTrace(trace, "funcAA") {
		t.Errorf("should contains funcAA")
	}

	if !haveTrace(trace, "funcBB") {
		t.Errorf("should contains funcBB")
	}
}
