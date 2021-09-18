package errors

import (
	"fmt"

	"github.com/payfazz/go-errors/v2/trace"
)

// see https://pkg.go.dev/fmt/#Errorf
func Errorf(format string, a ...interface{}) error {
	return &tracedErr{
		err:   fmt.Errorf(format, a...),
		trace: trace.Get(1, defaultDeep),
	}
}
