package errors_test

import (
	"fmt"
	"testing"

	"github.com/payfazz/go-errors"
)

type myErr struct{}

func (*myErr) Error() string {
	return ""
}

func Test0(t *testing.T) {
	var wrappedErr error
	wrappedErr = errors.Wrap(nil)
	if wrappedErr != nil {
		t.Errorf("Wrap(nil) should be nil")
	}
}
func Test1(t *testing.T) {
	var originalErr error
	originalErr = &myErr{}
	wrappedErr1 := errors.Wrap(originalErr)
	wrappedErr2 := errors.Wrap(wrappedErr1)
	wrappedErr3 := errors.Wrap(wrappedErr2)
	causeErr1 := errors.Cause(wrappedErr3)
	causeErr2 := errors.Cause(originalErr)

	if wrappedErr1 != wrappedErr2 || wrappedErr2 != wrappedErr3 {
		t.Errorf("wrapped error must be indempotent")
	}

	if causeErr1 != originalErr || causeErr2 != originalErr {
		t.Errorf("cause error mismatch")
	}
}

func Test2(t *testing.T) {
	var originalErr error
	originalErr = errors.New("testerr")
	wrappedErr1 := errors.Wrap(originalErr)
	wrappedErr2 := errors.Wrap(wrappedErr1)
	wrappedErr3 := errors.Wrap(wrappedErr2)
	causeErr1 := errors.Cause(wrappedErr3)
	causeErr2 := errors.Cause(originalErr)

	if wrappedErr1 != wrappedErr2 || wrappedErr2 != wrappedErr3 {
		t.Errorf("wrapped error must be indempotent")
	}

	if causeErr1 != originalErr || causeErr2 != originalErr || wrappedErr1 != originalErr {
		t.Errorf("cause error mismatch")
	}
}

func Test3(t *testing.T) {
	err := fmt.Errorf("testerr")
	wrappedErr := errors.Wrap(err)
	if wrappedErr.Error() != "testerr" {
		t.Errorf("Wrapped error should not change error message")
	}
}