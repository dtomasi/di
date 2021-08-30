package errors_test

import (
	"github.com/dtomasi/di/internal/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func funcThatReturnsAnInternalErr() error {
	return errors.NewError("unknown error")
}

func funcThatCallsOtherFuncAndGivesWrappedErrUp(path string) (err error) {
	defer errors.WrapErrf(&err, "funcThatReturnsAnInternalErr(%q)", path)

	err = funcThatReturnsAnInternalErr()
	if err != nil {
		return err
	}

	return nil
}

func TestWrapErr(t *testing.T) {
	err := funcThatCallsOtherFuncAndGivesWrappedErrUp("/not-existing")
	assert.Error(t, err)
	assert.IsType(t, (*errors.Error)(nil), err.(*errors.Error).Unwrap()) // nolint
}
