package errors

// PrintTo print formated error p
func PrintTo(p interface{ Print(...interface{}) }, err error) {
	if err == nil {
		return
	}
	if p == nil {
		return
	}
	p.Print(Format(Wrap(err)))
}
