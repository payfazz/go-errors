package errors_test

import (
	"strings"

	"github.com/payfazz/go-errors/v2/trace"
)

func funcAA(f func()) { f() }

func funcBB(f func()) { f() }

func haveTrace(ls []trace.Location, what string) bool {
	for _, l := range ls {
		if strings.Contains(l.Func(), what) {
			return true
		}
	}
	return false
}
