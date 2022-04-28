package di

import (
	"context"
	"fmt"
	loggerinternal "github.com/dtomasi/di/internal/pkg/logger"
	"github.com/dtomasi/di/internal/pkg/utils"
	"github.com/dtomasi/fakr"
	"github.com/dtomasi/go-event-bus/v3"
	z "github.com/dtomasi/zerrors"
	"github.com/go-logr/logr"
	"reflect"
)

// Container is the actual service container struct.
type Container struct {
	// execution context
	ctx context.Context

	// function to cancel the context
	ctxCancelFun context.CancelFunc

	// logger is used for internal logs.
	logger *loggerinternal.InternalLogger

	// eventBus is the eventbus instance
	eventBus *eventbus.EventBus

	// The ParameterProvider
	paramProvider ParameterProvider

	// Map of Service definitions
	serviceDefs *ServiceDefMap
}

// NewServiceContainer returns a new Container instance.
func NewServiceContainer(opts ...Option) *Container {
	c := &Container{ //nolint:exhaustivestruct
		ctx:           context.Background(),
		logger:        loggerinternal.NewInternalLogger(fakr.New()),
		eventBus:      eventbus.NewEventBus(),
		paramProvider: &NoParameterProvider{},
		serviceDefs:   NewServiceDefMap(),
	}

	for _, opt := range opts {
		opt(c)
	}

	// Get the context cancel function
	c.ctx, c.ctxCancelFun = context.WithCancel(c.ctx)

	// wrap container into context
	c.ctx = context.WithValue(c.ctx, ContextKeyContainer, c)

	return c
}

// SetParameterProvider allows to set the parameter provider even after container initialization.
func (c *Container) SetLogger(l logr.Logger) {
	// Wrap the original logger to apply module name
	c.logger = loggerinternal.NewInternalLogger(l)
}

// SetParameterProvider allows to set the parameter provider even after container initialization.
func (c *Container) SetParameterProvider(pp ParameterProvider) {
	c.paramProvider = pp
}

// GetEventBus returns the eventbus instance. This is used to register to internal events that can be used as hooks.
func (c *Container) GetEventBus() *eventbus.EventBus {
	return c.eventBus
}

// GetContext returns the context.
func (c *Container) GetContext() context.Context {
	return c.ctx
}

// GetContext returns the context.
func (c *Container) CancelContext() {
	c.ctxCancelFun()
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

	c.logger.Debug("added a new service via Set()", "service", ref.String())

	return c
}

// Get returns a requested service.
func (c *Container) Get(ref fmt.Stringer) (interface{}, error) {
	// s is a service object
	if sd, ok := c.serviceDefs.Load(ref); ok {
		if sd.instance == nil || sd.options.alwaysRebuild {
			var err error
			sd.instance, err = c.buildServiceInstance(sd)

			if err != nil {
				return nil, err
			}
		}

		return sd.instance, nil
	}

	return nil, z.NewWithOpts(
		fmt.Sprintf("services %s not found", ref),
		z.WithType(ServiceNotFoundError),
	)
}

// MustGet returns a service instance or panics on error.
func (c *Container) MustGet(ref fmt.Stringer) interface{} {
	i, err := c.Get(ref)
	if err != nil {
		panic(err)
	}

	return i
}

// FindByTags finds all service instances with given tags and returns them as a slice.
func (c *Container) FindByTags(tags []fmt.Stringer) ([]interface{}, error) {
	var instances []interface{}

	err := c.serviceDefs.Range(func(key fmt.Stringer, def *ServiceDef) error {
		matchCount := 0
		for _, searchTag := range tags {
			if containsTag(def.tags, searchTag) {
				matchCount++
			}
		}

		if matchCount == len(tags) {
			// use Get to ensure the service is built if not already.
			s, getErr := c.Get(key)
			if getErr != nil {
				return getErr
			}
			instances = append(instances, s)
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
	defer z.WrapPtrWithOpts(&err, "error while building container", z.WithType(ContainerBuildError))

	c.logger.Debug("starting container build")

	err = c.serviceDefs.Range(func(key fmt.Stringer, serviceDef *ServiceDef) error {
		c.logger.Debug("building services", "name", key.String())

		// skip lazy initializing services here
		if serviceDef.options.buildOnFirstRequest || serviceDef.options.alwaysRebuild {
			c.logger.Debug("skipping service because its set to lazy or should be rebuilt on each request",
				"name", key.String())

			return nil
		}

		// do not rebuild existing service instances
		if !serviceDef.options.alwaysRebuild && serviceDef.instance != nil {
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

	c.logger.Debug("container built successfully")
	c.eventBus.Publish(EventTopicDIReady.String(), c)

	return nil
}

func (c *Container) buildServiceInstance(def *ServiceDef) (instance interface{}, err error) {
	defer z.WrapPtrWithOpts(&err,
		fmt.Sprintf("error while building service %s", def.ref),
		z.WithType(ServiceBuildError),
	)

	parsedArgs, err := c.parseArgs(def)
	if err != nil {
		return nil, err
	}

	// build using provider function.
	if def.provider != nil {
		x := reflect.TypeOf(def.provider)

		if x.Kind() != reflect.Func {
			return nil, z.NewWithOpts("provider not a function", z.WithType(ProviderNotAFuncError))
		}

		inputArgCount := x.NumIn()
		if inputArgCount != len(parsedArgs) {
			return nil,
				z.NewWithOpts(fmt.Sprintf("expected %d got %d",
					inputArgCount,
					len(parsedArgs),
				), z.WithType(ProviderArgCountMismatchError))
		}

		var inputValues []reflect.Value

		for i := 0; i < inputArgCount; i++ {
			inType := x.In(i)

			inArgType := reflect.TypeOf(parsedArgs[i])

			if inArgType == nil {
				if inType.Kind() == reflect.Ptr {
					inputValues = append(inputValues, reflect.New(inType))
				} else {
					inputValues = append(inputValues, reflect.New(inType).Elem())
				}

				continue
			}

			inTypeString := utils.GetType(inType)
			inArgTypeString := utils.GetType(inArgType)

			if inTypeString != inArgTypeString && !inArgType.Implements(inType) {
				return nil,
					z.NewWithOpts(fmt.Sprintf("expected %s got %s",
						inTypeString,
						inArgTypeString,
					), z.WithType(ProviderArgTypeMismatchError))
			}

			inputValues = append(inputValues, reflect.ValueOf(parsedArgs[i]))
		}

		y := reflect.ValueOf(def.provider)
		returnValues := y.Call(inputValues)

		switch len(returnValues) {
		case 1:
			return returnValues[0].Interface(), nil
		case 2: // nolint:gomnd
			providerErr, ok := returnValues[1].Interface().(error)
			if !ok {
				providerErr = nil
			}

			return returnValues[0].Interface(), providerErr

		default:
			return nil,
				z.NewWithOpts(
					fmt.Sprintf("providers can only have 2 return values at max (interface{}, error). Got %d",
						len(returnValues),
					),
					z.WithType(ProviderToManyReturnValuesError),
				)
		}
	}

	return nil, z.NewWithOpts("provider missing", z.WithType(ProviderMissingError))
}

// parseArgs parses the arguments and assigns values by arg type.
// this function returns a new arg slice that is used for building the service
// without touching the original defined args.
func (c *Container) parseArgs(def *ServiceDef) (args []interface{}, err error) {
	c.logger.Debug("parsing args for provider of services", "name", def.ref.String())

	for _, v := range def.args {
		var val interface{}

		val, err = v.Evaluate(c)
		if err != nil {
			return
		}

		args = append(args, val)
	}

	c.logger.Debug("args for provider of services parsed successfully", "name", def.ref.String())

	return
}

func containsTag(a []fmt.Stringer, x fmt.Stringer) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}

	return false
}
