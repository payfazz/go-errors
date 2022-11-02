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

func TestCatchMultipleReturn(t *testing.T) {
	a, err := errors.Catch2(func() (int, error) { return 10, nil })
	if a != 10 || err != nil {
		t.FailNow()
	}

	a, b, err := errors.Catch3(func() (int, int, error) { return 10, 20, nil })
	if a != 10 || b != 20 || err != nil {
		t.FailNow()
	}

	a, b, c, err := errors.Catch4(func() (int, int, int, error) { return 10, 20, 30, nil })
	if a != 10 || b != 20 || c != 30 || err != nil {
		t.FailNow()
	}

	a, b, c, d, err := errors.Catch5(func() (int, int, int, int, error) { return 10, 20, 30, 40, nil })
	if a != 10 || b != 20 || c != 30 || d != 40 || err != nil {
		t.FailNow()
	}

	a, b, c, d, e, err := errors.Catch6(func() (int, int, int, int, int, error) { return 10, 20, 30, 40, 50, nil })
	if a != 10 || b != 20 || c != 30 || d != 40 || e != 50 || err != nil {
		t.FailNow()
	}

	orierr := errors.New("orierr")

	a, err = errors.Catch2(func() (int, error) { return 10, orierr })
	if a != 10 || !errors.Is(err, orierr) {
		t.FailNow()
	}

	a, b, err = errors.Catch3(func() (int, int, error) { return 10, 20, orierr })
	if a != 10 || b != 20 || !errors.Is(err, orierr) {
		t.FailNow()
	}

	a, b, c, err = errors.Catch4(func() (int, int, int, error) { return 10, 20, 30, orierr })
	if a != 10 || b != 20 || c != 30 || !errors.Is(err, orierr) {
		t.FailNow()
	}

	a, b, c, d, err = errors.Catch5(func() (int, int, int, int, error) { return 10, 20, 30, 40, orierr })
	if a != 10 || b != 20 || c != 30 || d != 40 || !errors.Is(err, orierr) {
		t.FailNow()
	}

	a, b, c, d, e, err = errors.Catch6(func() (int, int, int, int, int, error) { return 10, 20, 30, 40, 50, orierr })
	if a != 10 || b != 20 || c != 30 || d != 40 || e != 50 || !errors.Is(err, orierr) {
		t.FailNow()
	}
}
