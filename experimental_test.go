package errors_test

import (
	stderrors "errors"
	"testing"

	"github.com/payfazz/go-errors/v2"
)

type asTypeErrTest struct{}

func (*asTypeErrTest) Error() string { return "" }
func (*asTypeErrTest) SomeFunc()     {}

func TestAsType(t *testing.T) {
	originalErr := &asTypeErrTest{}
	wrappedErr := errors.Errorf("test %w", originalErr)

	e1, ok := errors.AsType[*asTypeErrTest](wrappedErr)
	if !ok {
		t.FailNow()
	}
	if e1 != originalErr {
		t.FailNow()
	}

	e2, ok := errors.AsType[interface {
		error
		SomeFunc()
	}](wrappedErr)
	if !ok {
		t.FailNow()
	}
	if e2 != originalErr {
		t.FailNow()
	}
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
