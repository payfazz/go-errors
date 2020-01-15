/*
Package errhandler provide utility to adding error handling.

With is the function to handle the error, this function must be called on deferred context.

for example:

	func main() {
		defer errhandler.With(nil)

		something, err := getSomething()
		errhandler.Check(errors.Wrap(err)) // using Wrap so we got the stack trace

		something2, err := getSomething2(something)
		if err != nil {
			errhandler.Check(errors.NewWithCause("getSomething2 is failing", err))
		}
	}

NOTE

please note that With adding some overhead, do not use it frequently, you should use golang idiom:

	if err != nil {
		return errors.Wrap(err)
	}

the only place to use this package is on main goroutine on main function
*/
package errhandler

import (
	"fmt"
	"os"

	"github.com/payfazz/go-errors"
)

type checkT struct{ error }

func (c checkT) Error() string {
	return "Unhandled error: If you got this message, it means that you forget to defer the error handler, " +
		"see: https://godoc.org/github.com/payfazz/go-errors/errhandler#With\n" +
		errors.Format(c.error)
}

// With will handle the error using f when Check is triggering the error
//
// if f is nil, default handler is to print error to stderr and exit with error code 1.
func With(f func(error)) {
	if rec := recover(); rec != nil {
		if c, ok := rec.(checkT); ok {
			if f == nil {
				f = defHandler
			}
			f(c.error)
		} else {
			panic(rec)
		}
	}
}

func defHandler(err error) {
	fmt.Fprint(os.Stderr, errors.Format(errors.Wrap(err)))
	os.Exit(1)
}

// Check the error
func Check(err error) {
	if err != nil {
		panic(checkT{err})
	}
}
