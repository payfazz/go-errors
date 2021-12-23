package errors_test

import (
	"fmt"
	"testing"

	"github.com/payfazz/go-errors/v2"
)

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
