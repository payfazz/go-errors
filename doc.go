// Package errors.
//
// This package can be used as drop-in replacement for standard errors package.
//
// This package provide StackTrace function to get the stack trace.
//
// Stack trace can be attached to any error by passing it to Trace function.
//
// New, NewWithCause, and Errorf function will return error that have stack trace.
package errors
