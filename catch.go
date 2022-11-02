package errors

import (
	stderrors "errors"
	"fmt"

	"github.com/payfazz/go-errors/v2/trace"
)

// run f, if f panic or returned, that value will be returned by this function
func Catch(f func() error) (err error) {
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

// like [Catch] but suitable for multiple return
func Catch2[A any](f func() (A, error)) (A, error) {
	var (
		a A
	)
	return a, Catch(func() error {
		var err error
		a, err = f()
		return err
	})
}

// like [Catch] but suitable for multiple return
func Catch3[A, B any](f func() (A, B, error)) (A, B, error) {
	var (
		a A
		b B
	)
	return a, b, Catch(func() error {
		var err error
		a, b, err = f()
		return err
	})
}

// like [Catch] but suitable for multiple return
func Catch4[A, B, C any](f func() (A, B, C, error)) (A, B, C, error) {
	var (
		a A
		b B
		c C
	)
	return a, b, c, Catch(func() error {
		var err error
		a, b, c, err = f()
		return err
	})
}

// like [Catch] but suitable for multiple return
func Catch5[A, B, C, D any](f func() (A, B, C, D, error)) (A, B, C, D, error) {
	var (
		a A
		b B
		c C
		d D
	)
	return a, b, c, d, Catch(func() error {
		var err error
		a, b, c, d, err = f()
		return err
	})
}

// like [Catch] but suitable for multiple return
func Catch6[A, B, C, D, E any](f func() (A, B, C, D, E, error)) (A, B, C, D, E, error) {
	var (
		a A
		b B
		c C
		d D
		e E
	)
	return a, b, c, d, e, Catch(func() error {
		var err error
		a, b, c, d, e, err = f()
		return err
	})
}
