package di

import z "github.com/dtomasi/zerrors"

// NoParameterProvider is a provider that is set by default.
type NoParameterProvider struct{}

func (p *NoParameterProvider) Get(_ string) (interface{}, error) {
	// Just return nil to not break call to Get if no parameter provider is set.
	return nil, z.NewWithOpts("provider not defined", z.WithType(ParamProviderNotDefinedError))
}
func (p *NoParameterProvider) Set(_ string, _ interface{}) error {
	// Same as above for the Setter here
	return z.NewWithOpts("provider not defined", z.WithType(ParamProviderNotDefinedError))
}
