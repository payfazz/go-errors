//go:build go1.20

package errors_test

import (
	"testing"

	"github.com/payfazz/go-errors/v2"
)

func TestJoin(t *testing.T) {
	a := errors.New("a")
	b := errors.New("b")
	c := errors.Join(a, b)

	errs := c.(interface{ Unwrap() []error }).Unwrap()
	if len(errs) != 2 {
		t.Fatalf("invalid Unwrap")
	}

	if !errors.Is(c, b) {
		t.Fatalf("invalid Is")
	}

	if !errors.Is(c, a) {
		t.Fatalf("invalid Is")
	}

	if errors.Unwrap(c) != nil {
		t.Fatalf("invalid unwrap")
	}

	if errors.Trace(errors.Trace(c)) != c {
		t.Fatalf("invalid Trace")
	}
}

func TestJoinErrorf(t *testing.T) {
	a := errors.New("a")
	b := errors.New("b")
	c := errors.Errorf("hai %w %w", a, b)

	errs := c.(interface{ Unwrap() []error }).Unwrap()
	if len(errs) != 2 {
		t.Fatalf("invalid Unwrap")
	}

	if !errors.Is(c, b) {
		t.Fatalf("invalid Is")
	}

	if !errors.Is(c, a) {
		t.Fatalf("invalid Is")
	}

	if errors.Unwrap(c) != nil {
		t.Fatalf("invalid unwrap")
	}

	if errors.Trace(errors.Trace(c)) != c {
		t.Fatalf("invalid Trace")
	}
}
