package errors

import (
	"strings"

	"github.com/payfazz/go-errors/v2/trace"
)

// like Format, but you can filter what location to include in the formated string
func FormatWithFilter(err error, filter func(trace.Location) bool) string {
	var sb strings.Builder

	add := func(err error) {
		if sb.Len() != 0 {
			sb.WriteString("Caused by ")
		}

		sb.WriteString("Error: ")
		sb.WriteString(err.Error())
		sb.WriteByte('\n')
		for _, l := range StackTrace(err) {
			if !filter(l) {
				continue
			}
			sb.WriteString("- ")
			sb.WriteString(l.String())
			sb.WriteByte('\n')
		}

		parentTrace := ParentStackTrace(err)
		firstParentTrace := true
		for _, l := range parentTrace {
			if !filter(l) {
				continue
			}
			if firstParentTrace {
				sb.WriteString("From goroutine created by:\n")
				firstParentTrace = false
			}
			sb.WriteString("- ")
			sb.WriteString(l.String())
			sb.WriteByte('\n')
		}
	}

	for err != nil {
		add(err)
		err = Unwrap(err)
	}

	return sb.String()
}

// like Format, but you can filter pkg location to include in the formated string
func FormatWithFilterPkgs(err error, pkgs ...string) string {
	return FormatWithFilter(err, func(l trace.Location) bool { return l.InPkg(pkgs...) })
}

// Format representation of the Error, including stack trace.
//
// Use err.Error() if you want to get just the error string
func Format(err error) string {
	return FormatWithFilter(err, func(l trace.Location) bool { return true })
}
