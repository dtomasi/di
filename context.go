package di

import (
	"context"
	"github.com/dtomasi/di/internal/errors"
)

const ContextKey = "di.container"

// GetContainerFromContext tries to get the container instance from given context as value.
func GetContainerFromContext(ctx context.Context) (*Container, error) {
	container, ok := ctx.Value(ContextKey).(*Container)
	if !ok {
		return container, errors.New("could not get container instance from context")
	}

	return container, nil
}
