package di

import (
	"context"
	"fmt"
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
	logger logr.Logger

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
		logger:        fakr.New(),
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

// SetLogger allows to pass a logr after container initialization.
func (c *Container) SetLogger(l logr.Logger) {
	c.logger = l
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

	c.logger.V(utils.LogLevelDebug).Info("added a new service via Set()", "service", ref.String())

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

	c.logger.V(utils.LogLevelDebug).Info("starting container build")

	err = c.serviceDefs.Range(func(key fmt.Stringer, serviceDef *ServiceDef) error {
		c.logger.V(utils.LogLevelDebug).Info("building services", "name", key.String())

		// skip lazy initializing services here
		if serviceDef.options.buildOnFirstRequest || serviceDef.options.alwaysRebuild {
			c.logger.V(utils.LogLevelDebug).Info("skipping service because its set to lazy or should be rebuilt on each request",
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

	c.logger.V(utils.LogLevelDebug).Info("container built successfully")
	c.eventBus.Publish(EventTopicDIReady.String(), c)

	return nil
}

func (c *Container) callReflectValueWithArgs(
	callable reflect.Value,
	serviceDefArgs []ServiceDefArg,
) (interface{}, error) {
	if callable.Type().Kind() != reflect.Func {
		return nil, z.NewWithOpts("callable not a function", z.WithType(CallableNotAFuncError))
	}

	evaluatedArgs, err := c.evaluateArgs(serviceDefArgs)
	if err != nil {
		return nil, err
	}

	callableNumInArgs := callable.Type().NumIn()
	if callableNumInArgs != len(evaluatedArgs) {
		return nil, z.NewWithOpts(fmt.Sprintf("expected %d got %d",
			callableNumInArgs,
			len(evaluatedArgs),
		), z.WithType(CallableArgCountMismatchError))
	}

	var callableInValues []reflect.Value

	// Prepare input args for callable
	for i := 0; i < callableNumInArgs; i++ {
		callableInType := callable.Type().In(i)

		inArgType := reflect.TypeOf(evaluatedArgs[i])

		if inArgType == nil {
			if callableInType.Kind() == reflect.Ptr {
				callableInValues = append(callableInValues, reflect.New(callableInType))
			} else {
				callableInValues = append(callableInValues, reflect.New(callableInType).Elem())
			}

			continue
		}

		inTypeString := utils.GetType(callableInType)
		inArgTypeString := utils.GetType(inArgType)

		if inTypeString != inArgTypeString && !inArgType.Implements(callableInType) {
			return nil,
				z.NewWithOpts(fmt.Sprintf("expected %s got %s",
					inTypeString,
					inArgTypeString,
				), z.WithType(CallableArgTypeMismatchError))
		}

		callableInValues = append(callableInValues, reflect.ValueOf(evaluatedArgs[i]))
	}

	// Call the callable
	returnValues := callable.Call(callableInValues)

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
				fmt.Sprintf("callable can only have 2 return values at max (interface{}, error). Got %d",
					len(returnValues),
				),
				z.WithType(CallableToManyReturnValuesError),
			)
	}
}

func (c *Container) buildServiceInstance(def *ServiceDef) (instance interface{}, err error) {
	defer z.WrapPtrWithOpts(&err,
		fmt.Sprintf("error while building service %s", def.ref),
		z.WithType(ServiceBuildError),
	)

	if def.provider == nil {
		return nil, z.NewWithOpts("provider missing", z.WithType(ProviderMissingError))
	}

	return c.callReflectValueWithArgs(reflect.ValueOf(def.provider), def.args)
}

// evaluateArgs parses the arguments and assigns values by arg type.
// this function returns a new arg slice that is used for building the service
// without touching the original defined args.
func (c *Container) evaluateArgs(args []ServiceDefArg) (eArgs []interface{}, err error) {
	for _, v := range args {
		var val interface{}

		val, err = v.Evaluate(c)
		if err != nil {
			return
		}

		eArgs = append(eArgs, val)
	}

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
