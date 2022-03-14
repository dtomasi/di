package di

import (
	"context"
	loggerinternal "github.com/dtomasi/di/internal/pkg/logger"
	eventbus "github.com/dtomasi/go-event-bus/v3"
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
func WithLogrImpl(l logr.Logger) Option {
	return func(c *Container) {
		c.logger = loggerinternal.NewInternalLogger(l)
	}
}

// WithEventBus defines an eventbus.EventBus instance to use instead of an internal one.
func WithEventBus(eb *eventbus.EventBus) Option {
	return func(c *Container) {
		c.eventBus = eb
	}
}
