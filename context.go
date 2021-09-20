package di

import (
	"context"
	z "github.com/dtomasi/zerrors"
)

type ContextKey int

const (
	ContextKeyContainer ContextKey = iota
)

// GetContainerFromContext tries to get the container instance from given context as value.
func GetContainerFromContext(ctx context.Context) (*Container, error) {
	container, ok := ctx.Value(ContextKeyContainer).(*Container)
	if !ok {
		return container, z.New("could not get container instance from context")
	}

	return container, nil
}
