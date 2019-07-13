package errors

import (
	"fmt"
	"runtime"
)

// DefaultDeep is the default deep when generating stack trace
var DefaultDeep = 20

// Location of execution, it cointain filename and linenumber
type Location struct {
	File string
	Line int
}

// String representation of Location
func (l Location) String() string {
	return fmt.Sprintf("%s:%d", l.File, l.Line)
}

func generateStackTrace(skip, max int) []Location {
	if max <= 0 {
		return nil
	}
	skip += 2
	if skip < 0 {
		skip = 0
	}

	ret := make([]Location, 0, max)
	ptrs := make([]uintptr, max)
	ptrsNum := runtime.Callers(skip, ptrs)
	if ptrsNum > 0 {
		frames := runtime.CallersFrames(ptrs)
		for {
			frame, more := frames.Next()
			if frame.File == "" {
				frame.File = "*unknown"
			}
			ret = append(ret, Location{
				File: frame.File,
				Line: frame.Line,
			})
			if !more {
				break
			}
		}
	}

	return ret
}
