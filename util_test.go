package errors_test

import (
	"strings"
	"testing"

	"github.com/payfazz/go-errors/v2"
)

func TestFormat(t *testing.T) {
	var err error
	funcAA(func() {
		funcBB(func() {
			err = errors.New("err1")
			err = errors.NewWithCause("err2", err)
		})
	})

	f := errors.Format(err)

	if !strings.Contains(f, "funcAA") ||
		!strings.Contains(f, "funcBB") ||
		!strings.Contains(f, "err1") ||
		!strings.Contains(f, "err2") {
		t.FailNow()
	}
}

func TestFormatFilter(t *testing.T) {
	var err error
	funcAA(func() {
		funcBB(func() {
			err = errors.New("err1")
			err = errors.NewWithCause("err2", err)
		})
	})

	f := errors.FormatWithFilterPkgs(err)

	if strings.Contains(f, "funcAA") ||
		strings.Contains(f, "funcBB") ||
		!strings.Contains(f, "err1") ||
		!strings.Contains(f, "err2") {
		t.FailNow()
	}
}
