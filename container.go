package di

import (
	"context"
	"fmt"
	"github.com/dtomasi/di/internal/errors"
	"github.com/dtomasi/di/internal/utils"
	"github.com/dtomasi/fakr"
	"github.com/dtomasi/go-event-bus/v2"
	"github.com/go-logr/logr"
	"reflect"
)

const (
	loggerName           string = "di"
	loggerVerbosityDebug int    = 6
	singleReturnValue    int    = 1
	doubleReturnValue    int    = 2
)

// Container is the actual service container struct.
type Container struct {
	ctx context.Context
	// injectableLogger is used for injection.
	injectableLogger logr.Logger
	// logger is used for internal logs.
	logger        logr.Logger
	eventBus      *eventbus.EventBus
	paramProvider ParameterProvider
	serviceDefs   *ServiceDefMap
}

// NewServiceContainer returns a new Container instance.
func NewServiceContainer(opts ...Option) *Container {
	i := &Container{ //nolint:exhaustivestruct
		ctx:              context.Background(),
		injectableLogger: fakr.New(),
		eventBus:         eventbus.NewEventBus(),
		paramProvider:    &NoParameterProvider{},
		serviceDefs:      NewServiceDefMap(),
	}

	for _, opt := range opts {
		opt(i)
	}

	// Setup logger name
	i.logger = i.injectableLogger.WithName(loggerName)

	// wrap container into context
	i.ctx = context.WithValue(i.ctx, ContextKey, i) // nolint:staticcheck

	return i
}

// GetEventBus returns the eventbus instance. This is used to register to internal events that can be used as hooks.
func (c *Container) GetEventBus() *eventbus.EventBus {
	return c.eventBus
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
func (c *Container) Get(ref fmt.Stringer) (service interface{}, err error) {
	defer errors.WrapPtrErrf(&err, "get service -> %s", ref.String())

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

	return nil, errors.NewStringerErr(ErrServiceNotFound)
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
func (c *Container) FindByTag(tag fmt.Stringer) (services []interface{}, err error) {
	defer errors.WrapPtrErrf(&err, "find service by tag -> %s", tag.String())

	var instances []interface{}

	err = c.serviceDefs.Range(func(key fmt.Stringer, def *ServiceDef) error {
		for _, defTag := range def.tags {
			if defTag == tag {
				// use Get to ensure the service is built if not already.
				s, getErr := c.Get(key)
				if getErr != nil {
					return getErr
				}
				instances = append(instances, s)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return instances, nil
}

// Build will build the service container.
func (c *Container) Build() (err error) {
	defer errors.WrapPtrErr(&err, "build container")

	c.debugLogger().Info("starting container build")

	err = c.serviceDefs.Range(func(key fmt.Stringer, serviceDef *ServiceDef) error {
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
		_, getErr := c.Get(key)
		if getErr != nil {
			return getErr
		}

		// return true as we want to get all build errors as an output here.
		return nil
	})

	if err != nil {
		return err
	}

	c.debugLogger().Info("container built successfully")
	c.eventBus.Publish(EventTopicDIReady.String(), c)

	return nil
}

func (c *Container) buildServiceInstance(def *ServiceDef) (instance interface{}, err error) {
	defer errors.WrapPtrErrf(&err, "build service ->  %s", def.ref)

	parsedArgs, err := c.parseArgs(def)
	if err != nil {
		return nil, err
	}

	// build using provider function.
	if def.provider != nil {
		x := reflect.TypeOf(def.provider)

		if x.Kind() != reflect.Func {
			return nil, errors.NewStringerErr(ErrProviderNotAFunc)
		}

		inputArgCount := x.NumIn()
		if inputArgCount != len(parsedArgs) {
			return nil, errors.WrapErrStringer(
				errors.NewErrf("expected %d got %d",
					inputArgCount,
					len(parsedArgs),
				),
				ErrProviderArgCountMismatch,
			)
		}

		var inputValues []reflect.Value

		for i := 0; i < inputArgCount; i++ {
			inType := x.In(i)

			inArgType := reflect.TypeOf(parsedArgs[i].value)
			inTypeString := utils.GetType(inType)
			inArgTypeString := utils.GetType(inArgType)

			if inTypeString != inArgTypeString && !inArgType.Implements(inType) {
				return nil,
					errors.WrapErrStringer(
						errors.NewErrf("expected %s got %s",
							inTypeString,
							inArgTypeString,
						),
						ErrProviderArgTypeMismatch,
					)
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
			return nil, errors.NewStringerErr(ErrProviderToManyReturnValues)
		}
	}

	return nil, errors.NewStringerErr(ErrProviderMissing)
}

// parseArgs parses the arguments and assigns values by arg type.
// this function returns a new arg slice that is used for building the service
// without touching the original defined args.
func (c *Container) parseArgs(def *ServiceDef) ([]Arg, error) {
	var parsedArgs []Arg

	c.debugLogger().Info("parsing args for provider of services", "name", def.ref.String())

	for _, v := range def.args {
		var argValue interface{}

		switch v._type {
		case ArgTypeService:
			s, err := c.Get(v.value.(fmt.Stringer))
			if err != nil {
				return nil, err
			}

			argValue = s
		case ArgTypeParam:
			val, err := c.paramProvider.Get(v.value.(string))
			if err != nil {
				return nil, errors.WrapErrStringer(
					errors.NewErrf("key: %s", v.value),
					ErrParamProviderGet,
				)
			}

			argValue = val
		case ArgTypeInterface:
			// Take the argument as it is
			argValue = v.value
		case ArgTypeContext:
			// Push the context
			argValue = c.ctx
		case ArgTypeLogger:
			// Push the logger
			argValue = c.injectableLogger
		case ArgTypeContainer:
			// Push the container itself
			argValue = c
		case ArgTypeParamProvider:
			// Push the parameter provider
			argValue = c.paramProvider
		}

		parsedArgs = append(parsedArgs, Arg{ //nolint:exhaustivestruct
			value: argValue,
		})
	}

	c.debugLogger().
		Info("args for provider of services parsed successfully", "name", def.ref.String())

	return parsedArgs, nil
}

// debugLogger returns the logger with debug verbosity.
func (c *Container) debugLogger() logr.Logger {
	return c.logger.V(loggerVerbosityDebug)
}
