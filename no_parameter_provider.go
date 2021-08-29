package di

// NoParameterProvider is a provider that is set by default and panics on use for ux reasons.
type NoParameterProvider struct{}

func (p *NoParameterProvider) Get(_ string) (interface{}, error) {
	panic("no parameter provider set")
}
func (p *NoParameterProvider) Set(_ string, _ interface{}) error {
	panic("no parameter provider set")
}
