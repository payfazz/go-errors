package errors

import (
	"github.com/payfazz/go-errors/v2/trace"
)

// run f, if f panic or returning non-nil error, that error will be returned
func Catch(f func() error) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			if recAsErr, ok := rec.(*tracedErr); ok {
				err = recAsErr
			} else {
				err = &tracedErr{error: &anyErr{data: rec}, trace: trace.Get(1, defaultDeep)}
			}
		}
	}()

	return f()
}
