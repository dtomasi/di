package di

import (
	"github.com/dtomasi/di/internal/errors"
)

// NoParameterProvider is a provider that is set by default.
type NoParameterProvider struct{}

func (p *NoParameterProvider) Get(key string) (interface{}, error) {
	// Just return nil to not break call to Get if no parameter provider is set.
	return nil, errors.WrapErrStringer(
		errors.NewErrf("key %s not found", key),
		ErrParamProviderNotDefined,
	)
}
func (p *NoParameterProvider) Set(key string, value interface{}) error {
	// Same as above for the Setter here
	return errors.WrapErrStringer(
		errors.NewErrf("cannot set key %s to value %v", key, value),
		ErrParamProviderNotDefined,
	)
}
