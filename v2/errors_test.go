package errors_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/payfazz/go-errors/v2"
)

func TestWrappingNil(t *testing.T) {
	err := errors.Wrap(nil)
	if err != nil {
		t.Errorf("Wrap(nil) should be nil")
	}
}

func TestIndempotentWrap(t *testing.T) {
	err1 := errors.New("testerr")
	err2 := errors.Wrap(err1)
	err3 := errors.Wrap(err2)
	err4 := errors.Wrap(err3)

	if err2 != err3 || err3 != err4 {
		t.Errorf("wrapped error must be indempotent")
	}
}

func TestWrapMessage(t *testing.T) {
	err1 := errors.Errorf("testerr")
	err2 := errors.Wrap(err1)
	if err2.Error() != "testerr" {
		t.Errorf("Wrapped error should not change error message")
	}
}

func TestNew(t *testing.T) {
	err0 := fmt.Errorf("err1")
	err1 := errors.Wrap(err0)
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

func funcAAAAAA(notpanic bool) error {
	err := errors.Wrap(&myErr{msg: "test err"})
	if notpanic {
		return err
	}
	panic(err)
}

var zero = 0

func funcBBBBBB(notpanic bool, f func() error) error {
	errCh := make(chan error)
	errors.Go(
		func(err error) {
			errCh <- err
			close(errCh)
		},
		func() error {
			if notpanic {
				return f()
			}

			// this will panic
			return fmt.Errorf("%d", 10/zero)
		},
	)
	return <-errCh
}

func TestFormat(t *testing.T) {
	funcName := "funcAAAAAA"
	err1 := funcBBBBBB(true, func() error { return funcAAAAAA(true) })

	haveTrace := false
	for _, t := range errors.StackTrace(err1) {
		if strings.Contains(t.Func(), funcName) {
			haveTrace = true
			break
		}
	}

	if !haveTrace {
		t.Errorf("invalid errors.StackTrace")
	}

	e2 := errors.NewWithCause("test cause", err1)

	if !strings.Contains(errors.Format(e2), funcName) {
		t.Errorf("invalid errors.Format")
	}
}

func TestErrorsAs(t *testing.T) {
	var target *myErr

	e := funcAAAAAA(true)

	if !errors.As(e, &target) {
		t.Errorf("invalid errors.As")
	}

	if e.Error() != target.msg {
		t.Errorf("invalid errors.As")
	}
}

func TestGo(t *testing.T) {
	err := funcBBBBBB(true, func() error { return fmt.Errorf("test err") })

	goodParent := false
	for _, l := range errors.ParentStackTrace(err) {
		if strings.Contains(l.Func(), "funcBBBBBB") {
			goodParent = true
			break
		}
	}

	if !goodParent {
		t.Errorf("invalid errors.ParentStackTrace")
	}
}
func TestGoTraced(t *testing.T) {
	err := funcBBBBBB(true, func() error { return funcAAAAAA(true) })

	goodTrace := false
	for _, l := range errors.StackTrace(err) {
		if strings.Contains(l.Func(), "funcAAAAAA") {
			goodTrace = true
			break
		}
	}

	if !goodTrace {
		t.Errorf("invalid errors.StackTrace")
	}

	goodParent := false
	for _, l := range errors.ParentStackTrace(err) {
		if strings.Contains(l.Func(), "funcBBBBBB") {
			goodParent = true
			break
		}
	}

	if !goodParent {
		t.Errorf("invalid errors.ParentStackTrace")
	}
}

func TestGoNil(t *testing.T) {
	err := funcBBBBBB(true, func() error { return nil })
	if err != nil {
		t.Errorf("invalid errors.Go")
	}
}

func TestGoPanic(t *testing.T) {
	err := funcBBBBBB(false, nil)

	goodTrace := false
	for _, l := range errors.StackTrace(err) {
		if strings.Contains(l.Func(), "funcBBBBBB") {
			goodTrace = true
			break
		}
	}

	if !goodTrace {
		t.Errorf("invalid errors.StackTrace")
	}

	goodParent := false
	for _, l := range errors.ParentStackTrace(err) {
		if strings.Contains(l.Func(), "funcBBBBBB") {
			goodParent = true
			break
		}
	}

	if !goodParent {
		t.Errorf("invalid errors.ParentStackTrace")
	}
}

func TestGoTracedPanic(t *testing.T) {
	err := funcBBBBBB(true, func() error { return funcAAAAAA(false) })

	goodTrace := false
	for _, l := range errors.StackTrace(err) {
		if strings.Contains(l.Func(), "funcAAAAAA") {
			goodTrace = true
			break
		}
	}

	if !goodTrace {
		t.Errorf("invalid errors.StackTrace")
	}

	goodParent := false
	for _, l := range errors.ParentStackTrace(err) {
		if strings.Contains(l.Func(), "funcBBBBBB") {
			goodParent = true
			break
		}
	}

	if !goodParent {
		t.Errorf("invalid errors.ParentStackTrace")
	}
}

func TestNilTrace(t *testing.T) {
	err := fmt.Errorf("test err")
	if len(errors.StackTrace(err)) != 0 {
		t.Errorf("invalid errors.StackTrace")
	}
	if len(errors.ParentStackTrace(err)) != 0 {
		t.Errorf("invalid errors.ParentStackTrace")
	}
}
