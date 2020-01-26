package errhandler_test

import (
	"testing"

	"github.com/payfazz/go-errors/errhandler"
)

type myErr struct{}

func (*myErr) Error() string {
	return ""
}

func Test1(t *testing.T) {
	fail := true

	defer func() {
		if fail {
			t.Fail()
		}
	}()

	var outErr error

	defer errhandler.With(func(err error) {
		if _, ok := err.(*myErr); ok {
			if err == outErr {
				fail = false
			}
		}
	})

	outErr = &myErr{}
	errhandler.Check(outErr)
}

func Test2(t *testing.T) {
	fail := true

	defer func() {
		if fail {
			t.Fail()
		}
	}()

	defer func() {
		if rec := recover(); rec != nil {
			if _, ok := rec.(*myErr); ok {
				fail = false
			}
		}
	}()

	defer errhandler.With(func(err error) {
		// do nothing
	})

	var err error
	err = &myErr{}
	panic(err)
}

func Test3(t *testing.T) {
	var err error
	testErr := &myErr{}

	func() {
		defer errhandler.With(errhandler.CatchAndSet(&err))
		errhandler.Check(testErr)
		t.Fatalf("This should be not called")
	}()

	if err != testErr {
		t.Fatalf("error should be same")
	}
}
