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
	if errors.ParentStackTrace(fmt.Errorf("testerror")) != nil {
		t.Errorf("errors.ParentStackTrace on non traced error should return nil")
	}
}

func TestCatch(t *testing.T) {
	err1 := errors.Catch(func() error { return nil })
	if err1 != nil {
		t.Errorf("errors.Catch should return nil when f returning nil")
	}

	check := func(f func() error) {
		err2 := errors.Catch(func() error {
			var err error
			funcAA(func() {
				err = f()
			})
			return err
		})
		if err2 == nil {
			t.Errorf("errors.Catch should return non-nil when f returning non-nil or panic")
		}
		if !haveTrace(errors.StackTrace(err2), "funcAA") {
			t.Errorf("errors.Catch trace should contains funcAA")
		}
	}

	check(func() error {
		return errors.New("testerr")
	})

	check(func() error {
		panic(errors.New("testerr"))
	})

	check(func() error {
		panic(fmt.Errorf("testerr"))
	})

	check(func() error {
		var something interface{ something() }
		// this trigger nil pointer exception
		something.something()
		return nil
	})

	check(func() error {
		panic("a test string")
	})
}

func doErrorsGo(f func() error) error {
	var err error
	funcAA(func() {
		doneCh := make(chan struct{})
		report := func(e error) {
			err = e
			close(doneCh)
		}
		errors.Go(report, func() error {
			var innerErr error
			funcBB(func() {
				innerErr = f()
			})
			return innerErr
		})
		<-doneCh
	})
	return err
}

func TestGo(t *testing.T) {
	check := func(shouldNil bool, shouldHaveStackTrace bool, f func() error) {
		err := doErrorsGo(f)

		if !shouldNil {
			if shouldHaveStackTrace {
				if !haveTrace(errors.StackTrace(err), "funcBB") {
					t.Errorf("errors.Go stack trace should contains funcBB")
				}
			}

			if !haveTrace(errors.ParentStackTrace(err), "funcAA") {
				t.Errorf("errors.Go stack trace should contains funcAA")
			}
		} else if err != nil {
			t.Errorf("errors.Go should report nil")
		}
	}

	check(false, false, func() error {
		return fmt.Errorf("testerr")
	})

	check(false, true, func() error {
		return errors.New("testerr")
	})

	check(true, false, func() error {
		return nil
	})
}

func TestFormat(t *testing.T) {
	// this test is just for code coverage, because errors.Format is not for machine
	// so the formated string is not standarized
	err := doErrorsGo(func() error {
		return errors.NewWithCause("testerr outer", errors.New("testerr inner"))
	})

	if !strings.Contains(errors.Format(err), "funcAA") {
		t.Errorf("errors.Format should contains funcAA")
	}

	if !strings.Contains(errors.Format(err), "funcBB") {
		t.Errorf("errors.Format should contains funcBB")
	}

	filter := func(trace.Location) bool {
		return false
	}

	if strings.Contains(errors.FormatWithFilter(err, filter), "funcAA") {
		t.Errorf("errors.FormatWithDeep should not contains funcAA")
	}

	if strings.Contains(errors.FormatWithFilter(err, filter), "funcBB") {
		t.Errorf("errors.FormatWithDeep should not contains funcBB")
	}
}
