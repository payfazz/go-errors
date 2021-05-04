package errors

type textErr struct {
	text  string
	cause error
}

func (e *textErr) Error() string {
	return e.text
}

func (e *textErr) Unwrap() error {
	return e.cause
}
