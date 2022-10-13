package errors_test

import (
	stderrors "errors"
	"strings"
	"testing"

	"github.com/payfazz/go-errors/v2"
	"github.com/payfazz/go-errors/v2/trace"
)

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

func TestCheck(t *testing.T) {
	err := errors.Catch(func() error {
		funcAA(func() {
			errors.Check(stderrors.New("testerr"))
		})
		return nil
	})

	if !haveTrace(errors.StackTrace(err), "funcAA") {
		t.Fatalf("should contain funcAA")
	}
}

func TestAssert(t *testing.T) {
	err := errors.Catch(func() error {
		funcAA(func() {
			errors.Assert(false)
		})
		return nil
	})

	if !haveTrace(errors.StackTrace(err), "funcAA") {
		t.Fatalf("should contain funcAA")
	}
}
