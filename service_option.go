package di

// ServiceOption defines an option function.
type ServiceOption func(so *serviceOptions)

// serviceDefOptions lifecycle options for definition.
type serviceOptions struct {
	buildOnFirstRequest bool
	alwaysRebuild       bool
}

// newServiceOptions returns a serviceOptions instance with defaults.
func newServiceOptions() *serviceOptions {
	return &serviceOptions{
		buildOnFirstRequest: false,
		alwaysRebuild:       false,
	}
}

// BuildOnFirstRequest option will not create an instance from on building the container.
// Instead, it will wait until the service is requested the first time.
func BuildOnFirstRequest() ServiceOption {
	return func(opts *serviceOptions) {
		opts.buildOnFirstRequest = true
	}
}

// BuildAlwaysRebuild defines that a service should be rebuilt on each request.
func BuildAlwaysRebuild() ServiceOption {
	return func(opts *serviceOptions) {
		opts.alwaysRebuild = true
	}
}
