package errors_test

import (
	"errors"
	interr "github.com/dtomasi/di/internal/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func funcThatReturnsAnInternalErr() error {
	return errors.New("unknown error") // nolint:goerr113
}

func funcThatCallsOtherFuncAndGivesWrappedErrUp(path string) (err error) {
	defer interr.WrapPtrErrf(&err, "funcThatReturnsAnInternalErr(%q)", path)

	err = funcThatReturnsAnInternalErr()
	if err != nil {
		return err
	}

	return nil
}

func TestWrapErr(t *testing.T) {
	err := funcThatCallsOtherFuncAndGivesWrappedErrUp("/not-existing")
	assert.Error(t, err)
}
