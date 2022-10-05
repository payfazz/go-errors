package errors

import (
	stderrors "errors"
	"fmt"

	"github.com/payfazz/go-errors/v2/trace"
)

// run f, if f panic or returned, that value will be returned by this function
func Catch(f func() error) error {
	_, err := Catch2(func() (struct{}, error) { return struct{}{}, f() })
	return err
}

// same with Catch, but with generic result type
func Catch2[Result any](f func() (Result, error)) (result Result, err error) {
	defer func() {
		rec := recover()
		if rec == nil {
			return
		}

		recErr, ok := rec.(error)
		if !ok {
			err = &traced{stderrors.New(fmt.Sprint(rec)), trace.Get(1, traceDeep)}
			return
		}

		cur := recErr
		for cur != nil {
			if _, ok := cur.(*traced); ok {
				err = recErr
				return
			}
			cur = stderrors.Unwrap(cur)
		}

		err = &traced{recErr, trace.Get(1, traceDeep)}
	}()

	return f()
}
