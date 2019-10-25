package errors

// Printer interface
type Printer interface {
	Print(...interface{})
}

// PrintTo print formated error to printer
func PrintTo(p Printer, err error) {
	if err == nil {
		return
	}
	p.Print(Format(Wrap(err)))
}

// Format the error as string
func Format(err error) string {
	if err == nil {
		return ""
	}

	if err2, ok := err.(*Error); ok {
		return err2.String()
	}

	return err.Error()
}
