package errors

import (
	"strings"
)

// Format representation of the Error, including stack trace.
//
// Use err.Error() if you want to get just the error string
func Format(err error) string {
	return FormatWithDeep(err, defaultDeep)
}

// like Format, but limit stack trace count by deep
func FormatWithDeep(err error, deep int) string {
	var sb strings.Builder

	add := func(err error) {
		if sb.Len() != 0 {
			sb.WriteString("Caused by ")
		}

		sb.WriteString("Error: ")
		sb.WriteString(err.Error())
		sb.WriteByte('\n')
		for i, l := range StackTrace(err) {
			if i == deep {
				break
			}
			sb.WriteString("- ")
			sb.WriteString(l.String())
			sb.WriteByte('\n')
		}

		parentTrace := ParentStackTrace(err)
		if len(parentTrace) > 0 {
			sb.WriteString("From goroutine created by:\n")
			for i, l := range parentTrace {
				if i == deep {
					break
				}
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
