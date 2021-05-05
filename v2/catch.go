package errors

import (
	"github.com/payfazz/go-errors/v2/trace"
)

// run f, if f panic or returned some error, that error will be returned by this function
func Catch(f func() error) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			if recAsTracedErr, ok := rec.(*tracedErr); ok {
				err = recAsTracedErr
			} else if recAsErr, ok := rec.(error); ok {
				err = &tracedErr{err: recAsErr, trace: trace.Get(1, defaultDeep)}
			} else {
				err = &tracedErr{err: &anyErr{data: rec}, trace: trace.Get(1, defaultDeep)}
			}
		}
	}()

	return f()
}
