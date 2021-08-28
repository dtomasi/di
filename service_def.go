package di

import (
	"fmt"
)

// StringRef defines a service reference as string.
// As ServiceDef implements fmt.Stringer as an interface for referencing services.
type StringRef string

// String implements fmt.Stringer interface method.
func (r StringRef) String() string {
	return string(r)
}

// ServiceDef is a definition of a service
// it describes how a service should be created and handled inside the
// service container.
type ServiceDef struct {
	ref      fmt.Stringer
	instance interface{}
	options  *serviceOptions
	provider interface{}
	args     []Arg
	tags     []fmt.Stringer
}

// NewServiceDef creates a new service definition.
func NewServiceDef(ref fmt.Stringer) *ServiceDef {
	i := &ServiceDef{
		ref:      ref,
		instance: nil,
		options:  newServiceOptions(),
		provider: nil,
		args:     []Arg{},
		tags:     []fmt.Stringer{},
	}

	return i
}

// Opts allows to set some options for lifecycle management.
func (sd *ServiceDef) Opts(opts ...ServiceOption) *ServiceDef {
	for _, opt := range opts {
		opt(sd.options)
	}

	return sd
}

// Provider defines a function that returns the actual serve instance.
// This function can also accept arguments that are described using the Args function.
func (sd *ServiceDef) Provider(provider interface{}) *ServiceDef {
	sd.provider = provider

	return sd
}

// Args accepts multiple constructor/provider function arguments.
func (sd *ServiceDef) Args(args ...Arg) *ServiceDef {
	sd.args = append(sd.args, args...)

	return sd
}

// AddTag allows to add a tag to a service definition.
func (sd *ServiceDef) AddTag(name fmt.Stringer) *ServiceDef {
	sd.tags = append(sd.tags, name)

	return sd
}
