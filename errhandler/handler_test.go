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

	defer errhandler.With(func(err error) {
		if _, ok := err.(*myErr); ok {
			fail = false
		}
	})

	var err error
	err = &myErr{}
	errhandler.Check(err)
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
