package errors

import (
	"github.com/payfazz/go-errors/v2/trace"
)

// Get paren goroutine stack trace of err
//
// return nil if err doesn't have stack trace
func ParentStackTrace(err error) []trace.Location {
	if e, ok := err.(*tracedErr); ok {
		return e.parent
	}
	return nil
}

// Spawn go routine
//
// run f, if f panic or returned, that value will be passed to report function
//
// the non-nil reported error will return non-nil when passed to ParentStackTrace
func Go(report func(error), f func() error) {
	parent := trace.Get(1, defaultDeep)

	doReport := func(err error) {
		if err == nil {
			report(nil)
			return
		}

		if t, ok := err.(*tracedErr); ok {
			t.parent = parent
		} else {
			err = &tracedErr{err: err, parent: parent}
		}

		report(err)
	}

	go func() { doReport(Catch(f)) }()
}
