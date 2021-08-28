package di

import (
	"context"
	"github.com/go-logr/logr"
)

// Option defines the option implementation.
type Option func(c *Container)

// WithContext allows to provide a context to di container.
func WithContext(ctx context.Context) Option {
	return func(c *Container) {
		c.ctx = ctx
	}
}

// WithParameterProvider defines an option to set ParameterProvider interface.
func WithParameterProvider(pp ParameterProvider) Option {
	return func(c *Container) {
		c.paramProvider = pp
	}
}

// WithLogrImpl defines the logr.Logger implementation to use
// For details see: https://github.com/go-logr/logr#a-minimal-logging-api-for-go
func WithLogrImpl(li logr.Logger) Option {
	return func(c *Container) {
		c.logger = li
	}
}
