package di

import (
	"fmt"
	"github.com/dtomasi/di/internal/errors"
)

// NoParameterProvider is a provider that is set by default.
type NoParameterProvider struct{}

func (p *NoParameterProvider) Get(key string) (interface{}, error) {
	// Just return nil to not break call to Get if no parameter provider is set.
	return nil, errors.NewStringerError(
		ErrParamProviderNotDefined,
		errors.WithDetail(fmt.Sprintf("key %s requested", key)),
	)
}
func (p *NoParameterProvider) Set(key string, value interface{}) error {
	// Same as above for the Setter here
	return errors.NewStringerError(
		ErrParamProviderNotDefined,
		errors.WithDetail(fmt.Sprintf("key %s value %v provided", key, value)),
	)
}
