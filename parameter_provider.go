package di

// ParameterProvider defines the interface which is used in Container internally while building a service instance
// The interface is kept simple to achieve flexibility. It should be easy to use koanf or viper behind the scenes for
// managing parameters passed to the container.
// By default, di.Container is initialized with a map[string]interface{} provider
type ParameterProvider interface {
	// Get defines how a parameter is fetched from a provider
	// we want to get an explicit error here to report it
	Get(key string) (interface{}, error)
	// Set allows to set additional parameters
	// if a value cannot be set we require an error to report it
	Set(key string, value interface{}) error
}
