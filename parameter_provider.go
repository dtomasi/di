package di

// ParameterProvider defines the interface which is used in Container internally while building a service instance
// The interface is kept simple to achieve flexibility. It should be easy to use koanf or viper behind the scenes for
// managing parameters passed to the container.
// By default, di.Container is initialized with a map[string]interface{} provider.
type ParameterProvider interface {
	// Get defines how a parameter is fetched from a provider.
	// We want to get an explicit error here to report it.
	Get(key string) (interface{}, error)
	// Set allows to set additional parameters.
	// If a value cannot be set we require an error to report it
	Set(key string, value interface{}) error
}

// NoParameterProvider is a provider that is set by default and panics on use for ux reasons.
type NoParameterProvider struct{}

func (p *NoParameterProvider) Get(_ string) (interface{}, error) {
	panic("no parameter provider set")
}
func (p *NoParameterProvider) Set(_ string, _ interface{}) error {
	panic("no parameter provider set")
}
