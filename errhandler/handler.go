/*
Package errhandler provide utility to adding error handling.

With is the function to handle the error, this function must be called on deferred context.

for example:

	func main() {
		defer errhandler.With(nil) // or defer errhandler.With(errhandler.Default)

		something, err := getSomething()
		errhandler.Check(errors.Wrap(err)) // using Wrap so we got the stack trace

		something2, err := getSomething2(something)
		if err != nil {
			errhandler.Fail(errors.NewWithCause("getSomething2 is failing", err))
		}
	}

NOTE

please note that With adding some overhead, do not use it frequently, you should use

	if err != nil {
		return errors.Wrap(err)
	}

the only place to use this package is on main function or the start of go routine
*/
package errhandler

import (
	"fmt"
	"os"

	"github.com/payfazz/go-errors"
)

type checkT struct {
	err error
}

func (c checkT) Error() string {
	return "Unhandled error: If you got this message, it means that you forget to defer the error handler, " +
		"see: github.com/payfazz/go-errors/errhandler\n" +
		errors.Format(c.err)
}

// With will handle the error using f when Check or Fail is triggering the error
//
// if f is nil Default is used.
func With(f func(error)) {
	if f == nil {
		f = Default
	}
	if rec := recover(); rec != nil {
		if c, ok := rec.(checkT); ok {
			f(c.err)
		} else {
			panic(rec)
		}
	}
}

// Default is the default error handler,
func Default(err error) {
	fmt.Fprint(os.Stderr, errors.Format(err))
	os.Exit(1)
}

// Check the error, if not nil, then trigger Fail
func Check(err error) {
	if err != nil {
		panic(checkT{errors_wrap(1, "", err, errors.DefaultDeep)})
	}
}

// Fail with the error, it assume that With is already deferred,
// to handle this error.
//
// DO NOT call Fail with err == nil
func Fail(err error) {
	panic(checkT{errors_wrap(1, "", err, errors.DefaultDeep)})
}
