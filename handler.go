package errors

type checkT struct {
	err error
}

func (c checkT) Error() string {
	return "Unhandled error: If you got this message, it means that you forget to defer the error handler, " +
		"see: github.com/payfazz/go-errors\n" +
		Format(c.err)
}

// Handle the error, and store it to *errptr, this function must be called on deferred context, i.e. `defer Handle(&err)`.
// please note that this function adding some overhead. You should use `if err != nil` idiom.
func Handle(errptr *error) {
	if rec := recover(); rec != nil {
		if c, ok := rec.(checkT); ok {
			*errptr = c.err
		} else {
			panic(rec)
		}
	}
}

// HandleWith do the same thing as Handle, but call f when error occurs
func HandleWith(f func(error)) {
	if rec := recover(); rec != nil {
		if c, ok := rec.(checkT); ok {
			f(c.err)
		} else {
			panic(rec)
		}
	}
}

// Check the error, will panic if err not nil, it assume that Handle is already deferred
// to handle this error.
func Check(err error) {
	if err != nil {
		panic(checkT{err})
	}
}
