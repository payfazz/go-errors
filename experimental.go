package errors

import stderrors "errors"

// same with As, but with generic result type.
//
// this is experimental
func AsType[E error](err error) (E, bool) {
	var target E
	ret := stderrors.As(err, &target)
	return target, ret
}

// Check will panic if err is not nil
//
// this is experimental
func Check(err error) {
	if err != nil {
		panic(Trace(err))
	}
}
