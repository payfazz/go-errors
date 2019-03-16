// Package errors is a utility to handle common error pattern in golang.
package errors

type checkT struct {
	err error
}

func (c checkT) Error() string {
	return "FATAL: If you got this message, it means that you forget to defer \"Handle\", see: github.com/payfazz/go-errors\n" +
		"    " + Format(c.err)
}

// Handle the error, and store it to *errptr, this function must be called on deferred context, i.e. `defer Handle(&err)`.
// please note that this function adding some overhead. You should use `if err != nil` idiom.
func Handle(errptr *error) {
	handler(recover(), func(err error) { *errptr = err })
}

// HandleWith do the same thing as Handle, but call f when error occurs
func HandleWith(f func(error)) {
	handler(recover(), f)
}

func handler(rec interface{}, f func(error)) {
	if rec == nil {
		return
	}
	if c, ok := rec.(checkT); ok {
		if f != nil {
			f(c.err)
		}
	} else {
		panic(rec)
	}
}

// Check the error, will panic if err not nil, it assume that Handle is already deferred.
func Check(err error) {
	MustNilOr(err, nil)
}

// MustNilOr do the same thing as Check, but call f before panic
func MustNilOr(err error, f func()) {
	if err != nil {
		if f != nil {
			f()
		}
		panic(checkT{err})
	}
}
