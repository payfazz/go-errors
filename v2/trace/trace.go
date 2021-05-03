// Package trace provide utility to get stack trace.
package trace

import (
	"fmt"
	"runtime"
)

// Location of execution
type Location struct {
	file     string
	line     int
	function string
}

// String representation of Location
func (l Location) String() string {
	if l.function == "" {
		return fmt.Sprintf("%s:%d", l.file, l.line)
	}

	return fmt.Sprintf("%s:%d (%s)", l.file, l.line, l.function)
}

// the File that this Location point to
func (l Location) File() string {
	return l.file
}

// the Line that this Location point to
func (l Location) Line() int {
	return l.line
}

// the Function that this location point to
func (l Location) Function() string {
	return l.function
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
	ptrs = ptrs[:ptrsNum]

	if ptrsNum > 0 {
		frames := runtime.CallersFrames(ptrs)
		for {
			frame, more := frames.Next()
			if frame.Line == 0 && frame.File == "" {
				continue
			}
			ret = append(ret, Location{
				function: frame.Function,
				file:     frame.File,
				line:     frame.Line,
			})
			if !more {
				break
			}
		}
	}

	return ret
}
