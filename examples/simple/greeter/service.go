package greeter

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
)

type Greeter struct {
	ctx        context.Context
	logger     logr.Logger
	salutation string
}

func NewGreeter(ctx context.Context, logger logr.Logger, salutation string) *Greeter {
	return &Greeter{
		ctx:        ctx,
		logger:     logger,
		salutation: salutation,
	}
}

func (g *Greeter) Greet(name string) {
	g.logger.V(6).Info("Greet() function called")
	fmt.Printf("%s, %s", g.salutation, name) //nolint:forbidigo
}
