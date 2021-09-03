package errors

import (
	"errors"
	"fmt"
)

// WrapErrStringer allows to pass a fmt.Stringer compatible variable. This makes it easy to use int Types in constants
// and pass them to go:generate stringer.
func WrapErrStringer(err error, t fmt.Stringer) error {
	return WrapErr(err, t.String())
}

// WrapErrf wraps the given err into a new error created from format string and args.
func WrapErrf(err error, format string, args ...interface{}) error {
	return WrapErr(err, fmt.Sprintf(format, args...))
}

// WrapErr wraps the given err into a new error created from given string.
func WrapErr(err error, msg string) error {
	return fmt.Errorf("%s: %w", msg, err)
}

// WrapPtrErrStringer allows to pass a fmt.Stringer compatible variable.
// This makes it easy to use int Types in constants
// and pass them to go:generate stringer.
func WrapPtrErrStringer(errp *error, t fmt.Stringer) {
	if *errp != nil {
		WrapPtrErr(errp, t.String())
	}
}

// WrapPtrErrf wraps the given err into a new error created from format string and args.
func WrapPtrErrf(errp *error, format string, args ...interface{}) {
	if *errp != nil {
		WrapPtrErr(errp, fmt.Sprintf(format, args...))
	}
}

// WrapPtrErr wraps the given err into a new error created from given string.
func WrapPtrErr(errp *error, msg string) {
	if *errp != nil {
		*errp = fmt.Errorf("%s: %w", msg, *errp)
	}
}

// NewStringerErr creates a new error from fmt.Stringer interface.
func NewStringerErr(s fmt.Stringer) error {
	return errors.New(s.String()) // nolint:goerr113
}

// NewErrf creates a new error from a format string and arguments.
func NewErrf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...) // nolint:goerr113
}

// New returns a new error from go builtin errors package.
func New(msg string) error {
	return errors.New(msg) // nolint:goerr113
}
