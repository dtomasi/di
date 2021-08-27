package di

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	singleReturnValue int = 1
	doubleReturnValue int = 2
)

var (
	// Container is a singleton/global instance.
	c *Container //nolint:gochecknoglobals

	errContainer        = errors.New("container error")
	errBuildingService  = errors.New("service build error")
	errParameterParsing = errors.New("param paring error")
)

// Container is the actual service container struct.
type Container struct {
	ctx           context.Context
	paramProvider ParameterProvider
	serviceDefs   *serviceMap
	services      *serviceMap
}

// NewServiceContainer returns a new Container instance.
func NewServiceContainer() *Container {
	return &Container{
		ctx:           nil,
		paramProvider: nil,
		serviceDefs:   newServiceMap(),
		services:      newServiceMap(),
	}
}

// DefaultContainer returns the default Container instance.
func DefaultContainer() *Container {
	if c == nil {
		c = NewServiceContainer()
	}

	return c
}

// SetParameterProvider allows to define a provider for fetching parameters using dot.notation.
func (c *Container) SetParameterProvider(provider ParameterProvider) {
	c.paramProvider = provider
}

// GetParameterProvider returns the set parameter provider.
func (c *Container) GetParameterProvider() ParameterProvider {
	return c.paramProvider
}

// Register lets you register a new ServiceDef to the container.
func (c *Container) Register(defs ...*ServiceDef) {
	for _, def := range defs {
		c.serviceDefs.Store(def.ref, def)
	}
}

// Set sets a service to container.
func (c *Container) Set(ref fmt.Stringer, service interface{}) *Container {
	c.services.Store(ref, service)

	return c
}

// Get returns a requested service.
func (c *Container) Get(ref fmt.Stringer) (interface{}, error) {
	if service, ok := c.services.Load(ref); ok {
		return service, nil
	}

	if def, ok := c.serviceDefs.Load(ref); ok { // build the service regardless of lazy.
		s, err := c.buildServiceFromDefinition(def.(*ServiceDef))
		if err != nil {
			return nil, err
		}

		// Save the created service instance to services map,
		c.services.Store(ref, s)

		return s, nil
	}

	return nil, newError(errContainer, fmt.Sprintf("service %s not found", ref))
}

// MustGet returns a service instance or panics on error.
func (c *Container) MustGet(ref fmt.Stringer) interface{} {
	i, err := c.Get(ref)
	if err != nil {
		panic(err)
	}

	return i
}

// Build will build the service container.
func (c *Container) Build(ctx context.Context) error {
	c.ctx = ctx

	var buildErrors []string

	c.serviceDefs.Range(func(key, value interface{}) bool {
		// skip lazy initializing services here
		if value.(*ServiceDef).options.buildOnFirstRequest {
			return true
		}

		_, err := c.Get(value.(*ServiceDef).ref)
		if err != nil {
			buildErrors = append(buildErrors, err.Error())
		}

		// return true as we want to get all build errors as an output here.
		return true
	})

	if len(buildErrors) > 0 {
		return newError(errContainer, fmt.Errorf(strings.Join(buildErrors, "\n"))) //nolint:goerr113
	}

	return nil
}

func (c *Container) buildServiceFromDefinition(def *ServiceDef) (interface{}, error) {
	parsedArgs, err := c.parseParameters(def)
	if err != nil {
		return nil, err
	}

	// build using provider function.
	if def.provider != nil {
		x := reflect.TypeOf(def.provider)

		if x.Kind() != reflect.Func {
			return nil, newError(
				errBuildingService,
				fmt.Sprintf("provider defined for service definition %s is not a function", def.ref),
			)
		}

		inputArgCount := x.NumIn()
		if inputArgCount != len(parsedArgs) {
			return nil, newError(errBuildingService, fmt.Sprintf(
				"expected %d arguments for %s provider. Got %d",
				inputArgCount,
				def.ref,
				len(parsedArgs),
			))
		}

		var inputValues []reflect.Value

		for i := 0; i < inputArgCount; i++ {
			inType := x.In(i)

			inArgType := reflect.TypeOf(parsedArgs[i].value)
			if getType(inType) != getType(inArgType) && !inArgType.Implements(inType) {
				return nil, newError(errBuildingService, fmt.Sprintf(
					"provider argument at position %d should be type of or implementing %s. Got %s",
					i,
					getType(inType),
					getType(inArgType),
				))
			}

			inputValues = append(inputValues, reflect.ValueOf(parsedArgs[i].value))
		}

		y := reflect.ValueOf(def.provider)
		returnValues := y.Call(inputValues)

		switch len(returnValues) {
		case singleReturnValue:
			return returnValues[0].Interface(), nil
		case doubleReturnValue:
			err, ok := returnValues[1].Interface().(error)
			if !ok {
				err = nil
			}

			return returnValues[0].Interface(), err

		default:
			return nil, newError(
				errBuildingService,
				fmt.Sprintf("to many return values in provider function for services %s", def.ref),
			)
		}
	}

	return nil, newError(errBuildingService, fmt.Sprintf("no provider function set for service %s", def.ref))
}

func (c *Container) parseParameters(def *ServiceDef) ([]Arg, error) {
	var parsedArgs []Arg

	for _, v := range def.args {
		switch v._type {
		case ArgTypeService:
			service, err := c.Get(v.value.(fmt.Stringer))
			if err != nil {
				return nil, newError(errParameterParsing, err)
			}

			parsedArgs = append(parsedArgs, Arg{ //nolint:exhaustivestruct
				value: service,
			})
		case ArgTypeParam:
			val, err := c.paramProvider.Get(v.value.(string))
			if err != nil {
				return nil, newError(errParameterParsing, err)
			}

			parsedArgs = append(parsedArgs, Arg{ //nolint:exhaustivestruct
				value: val,
			})
		case ArgTypeInterface:
			// Take the argument as it is
			parsedArgs = append(parsedArgs, Arg{ //nolint:exhaustivestruct
				value: v.value,
			})
		case ArgTypeContext:
			// Push the context
			parsedArgs = append(parsedArgs, Arg{ //nolint:exhaustivestruct
				value: c.ctx,
			})
		case ArgTypeContainer:
			// Push the container itself
			parsedArgs = append(parsedArgs, Arg{ //nolint:exhaustivestruct
				value: c,
			})
		case ArgTypeParamProvider:
			// Push the parameter provider
			parsedArgs = append(parsedArgs, Arg{ //nolint:exhaustivestruct
				value: c.paramProvider,
			})
		}
	}

	return parsedArgs, nil
}
