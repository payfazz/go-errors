package errors_test

import (
	"fmt"
	"testing"

	"github.com/payfazz/go-errors/v2"
	"github.com/payfazz/go-errors/v2/trace"
)

func TestTraceNil(t *testing.T) {
	err := errors.Trace(nil)
	if err != nil {
		t.Errorf("errors.Trace(nil) should be nil")
	}
}

func TestTraceMessage(t *testing.T) {
	err1 := errors.New("testerr")
	err2 := errors.Trace(err1)
	if err2.Error() != "testerr" {
		t.Errorf("errors.Trace should not change error message")
	}
}

func TestIndempotentTrace(t *testing.T) {
	err1 := errors.Errorf("testerr")
	err2 := errors.Trace(err1)
	err3 := errors.Trace(err2)

	if err1 != err2 || err2 != err3 {
		t.Errorf("traced error must be indempotent")
	}
}

func TestNewWithCause(t *testing.T) {
	err0 := fmt.Errorf("err1")
	err1 := errors.Trace(err0)
	err2 := errors.NewWithCause("err2", err1)
	err3 := errors.NewWithCause("err3", err2)

	if !errors.Is(err3, err2) {
		t.Errorf("invalid errors.Is")
	}

	if !errors.Is(err3, err1) {
		t.Errorf("invalid errors.Is")
	}

	if !errors.Is(err3, err0) {
		t.Errorf("invalid errors.Is")
	}
}

func TestTraceMessageErrorf(t *testing.T) {
	err1 := fmt.Errorf("testerr")
	err2 := errors.Errorf("testwrapper: %w", err1)
	if !errors.Is(errors.Unwrap(err2), err1) {
		t.Errorf("errors.Errorf should support %%w")
	}
}

func TestStackTrace(t *testing.T) {
	var err error
	funcAA(func() {
		funcBB(func() {
			err = errors.New("testerr")
		})
	})

	trace := errors.StackTrace(err)

	if !haveTrace(trace, "funcAA") {
		t.Errorf("errors.StackTrace should contains funcAA")
	}

	if !haveTrace(trace, "funcBB") {
		t.Errorf("errors.StackTrace should contains funcBB")
	}
}

func TestNonTraced(t *testing.T) {
	if errors.StackTrace(fmt.Errorf("testerror")) != nil {
		t.Errorf("errors.StackTrace on non traced error should return nil")
	}
}

type myErr struct{ msg string }

func (e *myErr) Error() string { return e.msg }

func TestErrorsAs(t *testing.T) {
	var target *myErr

	err := errors.Trace(&myErr{msg: "testerr"})

	if !errors.As(err, &target) {
		t.Errorf("invalid errors.As")
	}

	if err.Error() != target.msg {
		t.Errorf("invalid errors.As")
	}
}

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

func TestUnwrapInterface(t *testing.T) {
	ori := fmt.Errorf("test err")
	traced := errors.Trace(ori)

	untracealble, ok := traced.(interface{ Untrace() error })
	if !ok {
		t.Fatalf("should untraceable")
	}
	if untracealble.Untrace() != ori {
		t.Fatalf("should return same error")
	}
}
