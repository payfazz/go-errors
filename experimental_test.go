package errors

import (
	"testing"
)

type asTypeErrTest struct{}

func (*asTypeErrTest) Error() string { return "" }
func (*asTypeErrTest) SomeFunc()     {}

func TestAsType(t *testing.T) {
	originalErr := &asTypeErrTest{}
	wrappedErr := Errorf("test %w", originalErr)

	e1, ok := AsType[*asTypeErrTest](wrappedErr)
	if !ok {
		t.FailNow()
	}
	if e1 != originalErr {
		t.FailNow()
	}

	e2, ok := AsType[interface {
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
