package errors

import (
	"github.com/payfazz/go-errors/v2/trace"
)

// Get parent goroutine stack trace of err
//
// return nil if err doesn't have stack trace
//
// parent goroutine stack trace only available if the goroutine create by Go function
func ParentStackTrace(err error) []trace.Location {
	switch e := err.(type) {
	case *tracedErr:
		if e.parent == nil {
			return nil
		}
		return *e.parent
	case interface{ ParentStackTrace() []trace.Location }:
		return e.ParentStackTrace()
	default:
		return nil
	}
}

// Spawn go routine
//
// run f in that go routine, if f panic or returned, that value will be passed to report function,
//
// nil WILL be reported to report function if no error occured to indicate that f is finished
//
// the non-nil reported error will return non-nil when passed to ParentStackTrace
func Go(report func(error), f func() error) {
	parent := trace.Get(1, defaultDeep)

	doReport := func(data interface{}) {
		if data == nil {
			report(nil)
			return
		}

		if err, ok := data.(*tracedErr); ok {
			err.parent = &parent
			report(err)
			return
		}

		if err, ok := data.(error); ok {
			report(&tracedErr{
				err:    err,
				trace:  trace.Get(2, defaultDeep),
				parent: &parent,
			})
			return
		}

		report(&tracedErr{
			err:    &anyErr{data: data},
			trace:  trace.Get(2, defaultDeep),
			parent: &parent,
		})
	}

	go func() {
		defer func() {
			if rec := recover(); rec != nil {
				doReport(rec)
			}
		}()
		doReport(f())
	}()
}
