package di

import (
	"context"
	"errors"
	"fmt"
	"github.com/dtomasi/di/internal/utils"
	"github.com/dtomasi/fakr"
	"github.com/go-logr/logr"
	"reflect"
)

const (
	singleReturnValue int = 1
	doubleReturnValue int = 2
)

var (
	loggerName           = "di" //nolint:gochecknoglobals
	loggerVerbosityDebug = 6    //nolint:gochecknoglobals
	loggerVerbosityError = 1    //nolint:gochecknoglobals

	errContainer        = errors.New("container error")
	errBuildingService  = errors.New("service build error")
	errParameterParsing = errors.New("param paring error")
)

// Container is the actual service container struct.
type Container struct {
	ctx context.Context
	// originalLogger is used for injection.
	originalLogger logr.Logger
	// logger is used for internal logs.
	logger        logr.Logger
	paramProvider ParameterProvider
	serviceDefs   *ServiceDefMap
}

// NewServiceContainer returns a new Container instance.
func NewServiceContainer(opts ...Option) *Container {
	i := &Container{ //nolint:exhaustivestruct
		ctx:            context.Background(),
		originalLogger: fakr.New(),
		paramProvider:  &NoParameterProvider{},
		serviceDefs:    NewServiceDefMap(),
	}

	for _, opt := range opts {
		opt(i)
	}

	// Setup logger name
	i.logger = i.originalLogger.WithName(loggerName)

	return i
}

// Register lets you register a new ServiceDef to the container.
func (c *Container) Register(defs ...*ServiceDef) {
	for _, def := range defs {
		c.serviceDefs.Store(def.ref, def)
	}
}

// Set sets a service to container.
func (c *Container) Set(ref fmt.Stringer, s interface{}) *Container {
	c.serviceDefs.Store(ref, &ServiceDef{ //nolint:exhaustivestruct
		instance: s,
		options:  newServiceOptions(),
		tags:     []fmt.Stringer{},
	})

	c.debugLogger().Info("added a new service via Set()", "service", ref.String())

	return c
}

// Get returns a requested service.
func (c *Container) Get(ref fmt.Stringer) (interface{}, error) {
	var err error
	// s is a service object
	if sd, ok := c.serviceDefs.Load(ref); ok {
		if sd.instance == nil || sd.options.alwaysRebuild {
			sd.instance, err = c.buildServiceInstance(sd)
			if err != nil {
				return nil, err
			}
		}

		return sd.instance, nil
	}

	return nil, c.createAndLogError(errContainer, fmt.Errorf("service %s not found", ref).Error()) // nolint:goerr113
}

// MustGet returns a service instance or panics on error.
func (c *Container) MustGet(ref fmt.Stringer) interface{} {
	i, err := c.Get(ref)
	if err != nil {
		panic(err)
	}

	return i
}

// FindByTag finds all service instances with given tag and returns them as a slice.
func (c *Container) FindByTag(tag fmt.Stringer) ([]interface{}, error) {
	var instances []interface{}

	err := c.serviceDefs.Range(func(key fmt.Stringer, def *ServiceDef) error {
		for _, defTag := range def.tags {
			if defTag == tag {
				// use Get to ensure the service is built if not already.
				s, err := c.Get(key)
				if err != nil {
					return err
				}
				instances = append(instances, s)
			}
		}

		return nil
	})

	if err != nil {
		return nil, c.createAndLogError(errContainer, err)
	}

	return instances, nil
}

// Build will build the service container.
func (c *Container) Build() error {
	c.debugLogger().Info("starting container build")

	err := c.serviceDefs.Range(func(key fmt.Stringer, serviceDef *ServiceDef) error {
		c.debugLogger().Info("building services", "name", key.String())

		// skip lazy initializing services here
		if serviceDef.options.buildOnFirstRequest || serviceDef.options.alwaysRebuild {
			c.debugLogger().
				Info("skipping service because its set to lazy or should be rebuilt on each request",
					"name", key.String())

			return nil
		}

		// we just run get without expecting an instance is returned.
		// this will trigger build if definition instance is nil.
		_, err := c.Get(key)
		if err != nil {
			c.logger.V(loggerVerbosityError).Error(err, "creation of service failed", "name", key.String())

			return c.createAndLogError(errBuildingService, err)
		}

		// return true as we want to get all build errors as an output here.
		return nil
	})

	if err != nil {
		c.logger.V(loggerVerbosityError).Error(err, "container build failed")

		return c.createAndLogError(errBuildingService, err)
	}

	c.debugLogger().Info("container built successfully")

	return nil
}

func (c *Container) buildServiceInstance(def *ServiceDef) (interface{}, error) {
	parsedArgs, err := c.parseParameters(def)
	if err != nil {
		return nil, err
	}

	// build using provider function.
	if def.provider != nil {
		x := reflect.TypeOf(def.provider)

		if x.Kind() != reflect.Func {
			return nil, c.createAndLogError(
				errBuildingService,
				fmt.Sprintf("provider defined for service definition %s is not a function", def.ref),
			)
		}

		inputArgCount := x.NumIn()
		if inputArgCount != len(parsedArgs) {
			return nil, c.createAndLogError(errBuildingService, fmt.Sprintf(
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
			inTypeString := utils.GetType(inType)
			inArgTypeString := utils.GetType(inArgType)

			if inTypeString != inArgTypeString && !inArgType.Implements(inType) {
				return nil, c.createAndLogError(errBuildingService, fmt.Sprintf(
					"provider argument at position %d should be type of or implementing %s. Got %s",
					i,
					inTypeString,
					inArgTypeString,
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
			providerErr, ok := returnValues[1].Interface().(error)
			if !ok {
				providerErr = nil
			}

			return returnValues[0].Interface(), providerErr

		default:
			return nil, c.createAndLogError(
				errBuildingService,
				fmt.Sprintf("to many return values in provider function for services %s", def.ref),
			)
		}
	}

	return nil, c.createAndLogError(errBuildingService, fmt.Sprintf("no provider function set for service %s", def.ref))
}

// parseParameters parses the arguments and assigns values by arg type.
// this function returns a new arg slice that is used for building the service
// without touching the original defined args.
func (c *Container) parseParameters(def *ServiceDef) ([]Arg, error) {
	var parsedArgs []Arg

	c.debugLogger().Info("parsing parameters for provider of services", "name", def.ref.String())

	for _, v := range def.args {
		switch v._type {
		case ArgTypeService:
			s, err := c.Get(v.value.(fmt.Stringer))
			if err != nil {
				return nil, c.createAndLogError(errParameterParsing, err)
			}

			parsedArgs = append(parsedArgs, Arg{ //nolint:exhaustivestruct
				value: s,
			})
		case ArgTypeParam:
			val, err := c.paramProvider.Get(v.value.(string))
			if err != nil {
				return nil, c.createAndLogError(errParameterParsing, err)
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
		case ArgTypeLogger:
			// Push the logger
			parsedArgs = append(parsedArgs, Arg{ //nolint:exhaustivestruct
				value: c.logger,
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

	c.debugLogger().
		Info("parameters for provider of services parsed successfully", "name", def.ref.String())

	return parsedArgs, nil
}

// debugLogger returns the logger with debug verbosity.
func (c *Container) debugLogger() logr.Logger {
	return c.logger.V(loggerVerbosityDebug)
}

// createAndLogError logs and creates a new error that is returned.
func (c *Container) createAndLogError(errType error, msgOrErr interface{}) error { // Wrap string in error
	msgOrErr, ok := msgOrErr.(error)
	if !ok {
		msgOrErr = fmt.Errorf(msgOrErr.(string)) // nolint:goerr113
	}

	c.logger.V(loggerVerbosityError).Error(errType, msgOrErr.(error).Error())

	return fmt.Errorf("%w: %v", errType, msgOrErr)
}
