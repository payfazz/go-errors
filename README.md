# go-errors.

[![GoReference](https://pkg.go.dev/badge/github.com/payfazz/go-errors/v2)](https://pkg.go.dev/github.com/payfazz/go-errors/v2)

This package can be used as drop-in replacement for standard errors package.

This package provide `func StackTrace(error) []trace.Location` to get the stack trace.

Stack trace can be attached to any `error` by passing it to `func Trace(error) error`.

`New`, `NewWithCause`, and `Errorf` function will return error that have stack trace.
