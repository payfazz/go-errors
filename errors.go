package errors

import (
	stderrors "errors"
	"fmt"

	"github.com/payfazz/go-errors/v2/trace"
)

// see [stdlib errors.As]
//
// [stdlib errors.As]: https://pkg.go.dev/errors/#As
func As(err error, target any) bool {
	return stderrors.As(err, target)
}

// see [stdlib errors.Is]
//
// [stdlib errors.Is]: https://pkg.go.dev/errors/#Is
func Is(err, target error) bool {
	return stderrors.Is(err, target)
}

// see [stdlib errors.Unwrap]
//
// [stdlib errors.Unwrap]: https://pkg.go.dev/errors/#Unwrap
func Unwrap(err error) error {
	return stderrors.Unwrap(err)
}

// see [stdlib errors.New]
//
// [stdlib errors.New]: https://pkg.go.dev/errors/#New
func New(text string) error {
	return &traced{stderrors.New(text), trace.Get(1, traceDeep)}
}

// see [stdlib fmt.Errorf]
//
// [stdlib fmt.Errorf]: https://pkg.go.dev/fmt/#Errorf
func Errorf(format string, a ...any) error {
	return &traced{fmt.Errorf(format, a...), trace.Get(1, traceDeep)}
}

// will panic if err is not nil
func Check(err error) {
	if err != nil {
		panic(Trace(err))
	}
}

// Assert will panic if fact is false
func Assert(fact bool) {
	if !fact {
		panic(New("assertion failed"))
	}
}
