package errors_test

import (
	"testing"

	"github.com/payfazz/go-errors/v2"
	"github.com/payfazz/go-errors/v2/trace"
)

type myTracedError struct {
	trace  []trace.Location
	parent []trace.Location
}

func (m *myTracedError) Error() string                      { return "" }
func (m *myTracedError) StackTrace() []trace.Location       { return m.trace }
func (m *myTracedError) ParentStackTrace() []trace.Location { return m.parent }

func TestStackTraceInterface(t *testing.T) {
	var err error
	funcAA(func() {
		funcBB(func() {
			err = &myTracedError{trace.Get(0, 100), nil}
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

func TestGoInterface(t *testing.T) {
	var err error
	funcAA(func() {
		parent := trace.Get(0, 100)
		done := make(chan struct{})
		go func() {
			defer close(done)
			funcBB(func() {
				err = &myTracedError{trace.Get(0, 100), parent}
			})
		}()
		<-done
	})

	if !haveTrace(errors.StackTrace(err), "funcBB") {
		t.Errorf("stack trace should contains funcBB")
	}

	if !haveTrace(errors.ParentStackTrace(err), "funcAA") {
		t.Errorf("parent stack trace should contains funcAA")
	}

}
