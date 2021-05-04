package errors

import (
	"strings"
)

// Format representation of the Error, including stack trace.
//
// Use err.Error() if you want to get just the error string
func Format(err error) string {
	var sb strings.Builder

	add := func(err error) {
		if sb.Len() != 0 {
			sb.WriteString("Caused by ")
		}

		sb.WriteString("Error: ")
		sb.WriteString(err.Error())
		sb.WriteByte('\n')
		for _, l := range StackTrace(err) {
			sb.WriteString("- ")
			sb.WriteString(l.String())
			sb.WriteByte('\n')
		}

		parentTrace := ParentStackTrace(err)
		if len(parentTrace) > 0 {
			sb.WriteString("From goroutine created by:\n")
			for _, l := range parentTrace {
				sb.WriteString("- ")
				sb.WriteString(l.String())
				sb.WriteByte('\n')
			}
		}
	}

	for err != nil {
		add(err)
		err = Unwrap(err)
	}

	return sb.String()
}
