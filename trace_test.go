package errors_test

import (
	"fmt"
	"testing"

	"github.com/payfazz/go-errors/v2"
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

func TestErrorf(t *testing.T) {
	err0 := fmt.Errorf("err1")
	err1 := errors.Trace(err0)
	err2 := errors.Errorf("err2: %w", err1)
	err3 := errors.Errorf("err3: %w", err2)

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

func TestTraceMultipleReturn(t *testing.T) {
	a, err := errors.Trace2(10, nil)
	if a != 10 || err != nil {
		t.FailNow()
	}

	a, b, err := errors.Trace3(10, 20, nil)
	if a != 10 || b != 20 || err != nil {
		t.FailNow()
	}

	a, b, c, err := errors.Trace4(10, 20, 30, nil)
	if a != 10 || b != 20 || c != 30 || err != nil {
		t.FailNow()
	}

	orierr := errors.New("orierr")

	a, err = errors.Trace2(10, orierr)
	if a != 10 || !errors.Is(err, orierr) {
		t.FailNow()
	}

	a, b, err = errors.Trace3(10, 20, orierr)
	if a != 10 || b != 20 || !errors.Is(err, orierr) {
		t.FailNow()
	}

	a, b, c, err = errors.Trace4(10, 20, 30, orierr)
	if a != 10 || b != 20 || c != 30 || !errors.Is(err, orierr) {
		t.FailNow()
	}
}
