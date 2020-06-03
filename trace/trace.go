// Package trace provide utility to get stack trace.
package trace

import (
	"fmt"
	"runtime"
)

// Location of execution
type Location struct {
	File     string
	Line     int
	Function string
}

// String representation of Location
func (l Location) String() string {
	if l.Function == "" {
		return fmt.Sprintf("%s:%d", l.File, l.Line)
	}

	return fmt.Sprintf("%s:%d (%s)", l.File, l.Line, l.Function)
}

// Get return list of location of stack trace for calling function.
//
// max tell Get how deep the stack trace is.
// skip tell Get to skip some trace, 0 is where Get is called.
func Get(skip, max int) []Location {
	if max <= 0 {
		return nil
	}

	if skip < 0 {
		skip = 0
	}
	skip += 2

	ret := make([]Location, 0, max)

	ptrs := make([]uintptr, max)
	ptrsNum := runtime.Callers(skip, ptrs)
	if ptrsNum > 0 {
		frames := runtime.CallersFrames(ptrs)
		for {
			frame, more := frames.Next()
			if frame.Line == 0 && frame.File == "" {
				continue
			}
			ret = append(ret, Location{
				Function: frame.Function,
				File:     frame.File,
				Line:     frame.Line,
			})
			if !more {
				break
			}
		}
	}

	return ret
}
