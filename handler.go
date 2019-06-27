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

// Check the error, will panic if err not nil, it assume that Handle is already deferred.
func Check(err error) {
	if err != nil {
		panic(checkT{err})
	}
}

// MustNilOr do the same thing as Check, but call f before panic
func MustNilOr(err error, f func(error)) {
	if err != nil {
		f(err)
		panic(checkT{err})
	}
}

// CheckOrFail do the same thing as Fail, but only panic when err is not nil
func CheckOrFail(text string, err error) {
	if err != nil {
		panic(checkT{new(1, text, err, defaultDeep)})
	}
}

// CheckOrFailWithDeep is same with CheckOrFail, but with specified stack deep
func CheckOrFailWithDeep(text string, err error, deep int) {
	if err != nil {
		panic(checkT{new(1, text, err, deep)})
	}
}

// WrapAndCheck do the same thing as `Check(Wrap(err))`
func WrapAndCheck(err error) {
	if err != nil {
		panic(checkT{wrap(1, "", err, defaultDeep)})
	}
}

// WrapWithDeepAndCheck is same with WrapAndCheck, but with specified stack deep
func WrapWithDeepAndCheck(err error, deep int) {
	if err != nil {
		panic(checkT{wrap(1, "", err, deep)})
	}
}

// Fail do the same thing as `Check(NewWithCause(text, err))`.
func Fail(text string, err error) {
	panic(checkT{new(1, text, err, defaultDeep)})
}

// FailWithDeep is same with Fail, but with specified stack deep
func FailWithDeep(text string, err error, deep int) {
	panic(checkT{new(1, text, err, deep)})
}
