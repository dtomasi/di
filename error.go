package di

import "fmt"

// newError creates a new error containing the components prefix.
func newError(errType error, msg interface{}) error {
	return fmt.Errorf("%w: %v", errType, msg)
}
