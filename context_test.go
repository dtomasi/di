package di_test

import (
	"context"
	"github.com/dtomasi/di"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetContainerFromContext(t *testing.T) {
	c := di.NewServiceContainer()
	ctx := context.WithValue(context.Background(), di.ContextKeyContainer, c)

	ctxContainer, err := di.GetContainerFromContext(ctx)
	assert.NoError(t, err)
	assert.IsType(t, &di.Container{}, ctxContainer)
	assert.Equal(t, c, ctxContainer)
}
