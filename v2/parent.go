package errors

import (
	"github.com/payfazz/go-errors/v2/trace"
)

// Get stack trace of parent go routine of the error, return nil if none
func ParentStackTrace(err error) []trace.Location {
	if e, ok := err.(*tracedErr); ok {
		return e.parent
	}
	return nil
}

// Spawn new go routine
//
// run f, if f panic or returning non-nil error, it will be passed to the error channel
//
// the error passed to the error chanel will return non-nil when passed to ParentStackTrace
func Spawn(f func() error) <-chan error {
	errCh := make(chan error, 1)
	parent := trace.Get(1, defaultDeep)
	go func() {
		passErr := func(err error) {
			if err == nil {
				return
			}

			if t, ok := err.(*tracedErr); ok {
				t.parent = parent
				errCh <- t
				return
			}

			errCh <- &tracedErr{
				error:  err,
				parent: parent,
			}
		}

		defer func() {
			if rec := recover(); rec != nil {
				var err error
				if recAsErr, ok := rec.(*tracedErr); ok {
					err = recAsErr
				} else {
					err = &tracedErr{
						error: &anyErr{data: rec},
						trace: trace.Get(1, defaultDeep),
					}
				}
				passErr(err)
			}
			close(errCh)
		}()

		passErr(f())
	}()
	return errCh
}
