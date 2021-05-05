# go-errors.

[![GoDoc](https://pkg.go.dev/badge/github.com/payfazz/go-errors/v2)](https://pkg.go.dev/github.com/payfazz/go-errors/v2)

This package can be used as drop-in replacement for https://golang.org/pkg/errors package.

This package provide `StackTrace` function to get the stack trace.

Stack trace can be attached to any `error` by passing it to `Trace` function.

`New`, `NewWithCause`, and `Errorf` function will return error that have stack trace.
