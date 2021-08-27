package di

import (
	"context"
	"fmt"
)

// Container is the actual service container struct
type Container struct {
	ctx context.Context
	paramProvider ParameterProvider
	serviceDefs   *serviceMap
	services      *serviceMap
}

// Container is a singleton/global instance
var c *Container

// NewServiceContainer returns a new Container instance
func NewServiceContainer() *Container {
	return &Container{
		serviceDefs: newServiceMap(),
		services: newServiceMap(),
	}
}

// DefaultContainer returns the default Container instance
func DefaultContainer() *Container {
	if c == nil {
		c = NewServiceContainer()
	}
	return c
}

// SetParameterProvider allows to define a provider for fetching parameters using dot.notation
func (c *Container) SetParameterProvider(provider ParameterProvider) {
	c.paramProvider = provider
}

// GetParameterProvider returns the set parameter provider
func (c *Container) GetParameterProvider() ParameterProvider {
	return c.paramProvider
}

// Register lets you register a new ServiceDef to the container
func (c *Container) Register(defs ...*ServiceDef) {
	for _, def := range defs {
		c.serviceDefs.Store(def.ref, def)
	}
}

// Set sets a service to container
func (c *Container) Set(ref fmt.Stringer, service interface{}) *Container {
	c.services.Store(ref, service)
	return c
}

// Get returns a requested service
func (c *Container) Get(ref fmt.Stringer) (interface{}, error) {

	if service, ok := c.services.Load(ref); ok {
		return service, nil
	}

	if def, ok := c.serviceDefs.Load(ref); ok {

		// build the service regardless of lazy
		s, err := c.buildService(def.(*ServiceDef))
		if err != nil {
			return nil, err
		}

		// Save the created service instance to services map
		c.services.Store(ref, s)
		return s, nil
	}

	return nil, fmt.Errorf("service %s not found", ref)
}

// MustGet returns a service instance or panics on error
func (c *Container) MustGet(ref fmt.Stringer) interface{} {
	i, err := c.Get(ref)
	if err != nil {
		panic(err)
	}
	return i
}

// Build will build the service container
func (c *Container) Build(ctx context.Context) error {
	c.ctx = ctx

	var err error
	c.serviceDefs.Range(func(key, value interface{}) bool {
		// skip lazy initializing services here
		if value.(*ServiceDef).options.buildOnFirstRequest {
			return true
		}

		_, err = c.Get(value.(*ServiceDef).ref)
		return err == nil
	})

	return err
}

func (c *Container) buildService(def *ServiceDef) (interface{}, error) {

	// create a new arguments slice here used for def.build. We want to keep the original definition
	var parsedArgs []Arg
	for _, v := range def.args {
		switch v._type {
		case ArgTypeService:
			service, err := c.Get(v.value.(fmt.Stringer))
			if err != nil {
				return nil, fmt.Errorf("service get error: %v", err)
			}
			parsedArgs = append(parsedArgs, Arg{
				value: service,
			})
		case ArgTypeParam:
			val, err := c.paramProvider.Get(v.value.(string))
			if err != nil {
				return nil, fmt.Errorf("parameter provider error: %v", err)
			}
			parsedArgs = append(parsedArgs, Arg{
				value: val,
			})
		case ArgTypeInterface:
			// Take the argument as it is
			parsedArgs = append(parsedArgs, Arg{
				value: v.value,
			})
		case ArgTypeContext:
			// Push the context
			parsedArgs = append(parsedArgs, Arg{
				value: c.ctx,
			})
		case ArgTypeContainer:
			// Push the container itself
			parsedArgs = append(parsedArgs, Arg{
				value: c,
			})
		case ArgTypeParamProvider:
			// Push the parameter provider
			parsedArgs = append(parsedArgs, Arg{
				value: c.paramProvider,
			})
		}
	}

	i, err := def.build(parsedArgs)
	if err != nil {
		return nil, err
	}
	return i, nil
}
