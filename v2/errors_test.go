package errors_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/payfazz/go-errors/v2"
)

func TestWrappingNil(t *testing.T) {
	wrappedErr := errors.Wrap(nil)
	if wrappedErr != nil {
		t.Errorf("Wrap(nil) should be nil")
	}
}

func TestIndempotentWrap(t *testing.T) {
	errOri := errors.Errorf("testerr")
	errWrapped1 := errors.Wrap(errOri)
	errWrapped2 := errors.WrapWithDeep(20, errWrapped1)
	errWrapped3 := errors.Wrap(errWrapped2)

	if errWrapped1 != errWrapped2 || errWrapped2 != errWrapped3 {
		t.Errorf("wrapped error must be indempotent")
	}
}

func TestWrapMessage(t *testing.T) {
	errOri := fmt.Errorf("testerr")
	errWrapped := errors.Wrap(errOri)
	if errWrapped.Error() != "testerr" {
		t.Errorf("Wrapped error should not change error message")
	}
}

func TestNew(t *testing.T) {
	err0 := fmt.Errorf("err1")
	err1 := errors.Wrap(err0)
	err2 := errors.NewWithCause("err2", err1)
	err3 := errors.NewWithCause("err3", err2)
	err4 := errors.NewWithCause("err4", err3)
	err5 := errors.NewWithCause("err5", err4)

	if !errors.Is(err5, err0) {
		t.Errorf("invalid errors.Is")
	}

	if !errors.Is(err5, err4) {
		t.Errorf("invalid errors.Is")
	}

	if !errors.Is(err5, err3) {
		t.Errorf("invalid errors.Is")
	}

	if !errors.Is(err5, err2) {
		t.Errorf("invalid errors.Is")
	}

	if !errors.Is(err5, err1) {
		t.Errorf("invalid errors.Is")
	}
}

type myErr struct{ msg string }

func (e *myErr) Error() string { return e.msg }

// this function name is used for check stack trace
func func201df1b9f41c6cbf81ee12d34e90e26b() error {
	return errors.Wrap(&myErr{msg: "test err"})
}

func TestStackTrace(t *testing.T) {
	funcName := "func201df1b9f41c6cbf81ee12d34e90e26b"
	e1 := func201df1b9f41c6cbf81ee12d34e90e26b()

	haveTrace := false
	for _, t := range errors.StackTrace(e1) {
		if strings.Contains(t.Func(), funcName) {
			haveTrace = true
			break
		}
	}

	if !haveTrace {
		t.Errorf("invalid errors.StackTrace")
	}

	e2 := errors.NewWithCause("test cause", e1)

	if !strings.Contains(errors.Format(e2), funcName) {
		t.Errorf("invalid errors.Format")
	}
}

func TestErrorsAs(t *testing.T) {
	var target *myErr

	e := func201df1b9f41c6cbf81ee12d34e90e26b()

	if !errors.As(e, &target) {
		t.Errorf("invalid errors.As")
	}

	if e.Error() != target.msg {
		t.Errorf("invalid errors.As")
	}
}
