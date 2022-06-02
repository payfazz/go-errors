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
