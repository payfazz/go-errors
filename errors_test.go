package errors_test

import (
	"fmt"
	"strings"
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

func TestTraceMessageErrorf(t *testing.T) {
	err1 := fmt.Errorf("testerr")
	err2 := errors.Errorf("testwrapper: %w", err1)
	if !errors.Is(errors.Unwrap(err2), err1) {
		t.Errorf("errors.Errorf should support %%w")
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

func funcAA(f func()) { f() }

func funcBB(f func()) { f() }

func haveTrace(ls []trace.Location, what string) bool {
	for _, l := range ls {
		if strings.Contains(l.Func(), what) {
			return true
		}
	}
	return false
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

func TestCatch(t *testing.T) {
	check := func(shouldNil bool, shouldHaveTrace bool, f func() error) {
		err := errors.Catch(func() error {
			var err error
			funcAA(func() {
				err = f()
			})
			return err
		})

		if shouldNil {
			if err != nil {
				t.Errorf("errors.Catch should return nil when f returning nil")
			}
		} else if err == nil {
			t.Errorf("errors.Catch should return non-nil when f returning non-nil or panic")
		}

		if shouldHaveTrace && !haveTrace(errors.StackTrace(err), "funcAA") {
			t.Errorf("errors.Catch trace should contains funcAA")
		}
	}

	check(true, false, func() error {
		return nil
	})

	check(false, true, func() error {
		return errors.New("testerr")
	})

	check(false, false, func() error {
		return fmt.Errorf("testerr")
	})

	check(false, true, func() error {
		panic(errors.New("testerr"))
	})

	check(false, true, func() error {
		panic(fmt.Errorf("testerr"))
	})

	check(false, true, func() error {
		var something interface{ something() }
		// this trigger nil pointer exception
		something.something()
		return nil
	})

	check(false, true, func() error {
		panic("a test string")
	})
}
