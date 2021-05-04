package errors

import "fmt"

type anyErr struct {
	data  interface{}
	cause error
}

func (e *anyErr) Error() string {
	return fmt.Sprint(e.data)
}

func (e *anyErr) Unwrap() error {
	return e.cause
}
