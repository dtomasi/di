package errors

import (
	"fmt"
	"io"
	"strconv"
)

// WrapErrType allows to pass a fmt.Stringer compatible variable. This makes it easy to use int Types in constants
// and pass them to go:generate stringer.
func WrapErrType(errp *error, t fmt.Stringer) {
	if *errp != nil {
		WrapErr(errp, t.String())
	}
}

// WrapErrf wraps the given err into a new error created from format string and args.
func WrapErrf(errp *error, format string, args ...interface{}) {
	if *errp != nil {
		WrapErr(errp, fmt.Sprintf(format, args...))
	}
}

// WrapErr wraps the given err into a new error created from given string.
func WrapErr(errp *error, msg string) {
	if *errp != nil {
		WrapErrWithOpts(errp, msg)
	}
}

// WrapErrWithOpts same like WrapErr but with options.
// see: ErrOption.
func WrapErrWithOpts(errp *error, msg string, opts ...ErrOption) {
	if *errp != nil {
		opts = append(opts, WithPreviousErr(*errp))
		*errp = NewError(msg, opts...)
	}
}

// Option defines the option implementation.
type ErrOption func(c *Error)

// WithDetail allows to provide error details.
func WithDetail(detail string) ErrOption {
	return func(e *Error) {
		e.detail = detail
	}
}

// WithDetail allows to provide error details.
func WithPreviousErr(err error) ErrOption {
	return func(e *Error) {
		e.err = err
	}
}

// code from: https://github.com/jba/errfmt

// Error represent an error type that holds additional information in detail property.
type Error struct {
	msg, detail string
	err         error
}

// NewError creates a new Error
// It allows passing options like adding error details or previous err to wrap.
func NewError(msg string, opts ...ErrOption) *Error {
	err := &Error{
		msg:    msg,
		detail: "",
		err:    nil,
	}

	for _, opt := range opts {
		opt(err)
	}

	return err
}

// NewSError creates a new Error like NewError, but returns only error interface/type.
func NewStringerError(t fmt.Stringer, opts ...ErrOption) *Error {
	return NewError(t.String(), opts...)
}

// Unwrap returns the original (wrapped) error.
func (e *Error) Unwrap() error { return e.err }

// Error returns the error msg as string.
func (e *Error) Error() string {
	if e.err == nil {
		return e.msg
	}

	return e.msg + ": " + e.err.Error()
}

// Format pretty prints error and details to fmt.
func (e *Error) Format(s fmt.State, c rune) {
	if s.Flag('#') && c == 'v' {
		type nomethod Error

		_, _ = fmt.Fprintf(s, "%#v", (*nomethod)(e))

		return
	}

	if !s.Flag('+') || c != 'v' {
		_, _ = fmt.Fprintf(s, spec(s, c), e.Error())

		return
	}

	_, _ = fmt.Fprintln(s, e.msg)

	if e.detail != "" {
		_, _ = io.WriteString(s, "\t")
		_, _ = fmt.Fprintln(s, e.detail)
	}

	if e.err != nil {
		if ferr, ok := e.err.(fmt.Formatter); ok { // nolint:errorlint
			ferr.Format(s, c)
		} else {
			_, _ = fmt.Fprintf(s, spec(s, c), e.err)
			_, _ = io.WriteString(s, "\n")
		}
	}
}

func spec(s fmt.State, c rune) string {
	buf := []byte{'%'}

	for _, f := range []int{'+', '-', '#', ' ', '0'} {
		if s.Flag(f) {
			buf = append(buf, byte(f))
		}
	}

	if w, ok := s.Width(); ok {
		buf = strconv.AppendInt(buf, int64(w), 10) //nolint:gomnd
	}

	if p, ok := s.Precision(); ok {
		buf = append(buf, '.')
		buf = strconv.AppendInt(buf, int64(p), 10) //nolint:gomnd
	}

	buf = append(buf, byte(c))

	return string(buf)
}
