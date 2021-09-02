package errors_test

import (
	"github.com/dtomasi/di/internal/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func funcThatReturnsAnInternalErr() error {
	return errors.NewErrf("unknown error")
}

func funcThatCallsOtherFuncAndGivesWrappedPtrErrUp(path string) (err error) {
	defer errors.WrapPtrErrf(&err, "funcThatReturnsAnInternalErr(%q)", path)

	err = funcThatReturnsAnInternalErr()
	if err != nil {
		return err
	}

	return nil
}

func funcThatCallsOtherFuncAndGivesWrappedErrUp(path string) error {
	err := funcThatReturnsAnInternalErr()
	if err != nil {
		return errors.WrapErrf(err, "funcThatReturnsAnInternalErr(%q)", path)
	}

	return nil
}

func TestWrapPtrErr(t *testing.T) {
	err := funcThatCallsOtherFuncAndGivesWrappedPtrErrUp("/not-existing")
	assert.Error(t, err)
}

func TestWrapErr(t *testing.T) {
	err := funcThatCallsOtherFuncAndGivesWrappedErrUp("/not-existing")
	assert.Error(t, err)
}
