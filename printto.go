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
