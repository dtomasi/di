package di

import (
	"fmt"
	"reflect"
)

// ServiceDefArg is the interface that arguments have to im.
type ServiceDefArg interface {
	evaluate(*Container) (interface{}, error)
}

type interfaceArg struct {
	inValue interface{}
}

func (a *interfaceArg) evaluate(_ *Container) (interface{}, error) {
	return a.inValue, nil
}

// InterfaceArg is a shortcut for an argument of type interface{}.
// This argument allows to pass any value that is provided to the service.
func InterfaceArg(in interface{}) ServiceDefArg {
	return &interfaceArg{inValue: in}
}

// ServiceRef Arg
// This argument type allows to pass a reference to a service that will be injected.
type serviceRefArg struct {
	ref fmt.Stringer
}

func (a *serviceRefArg) evaluate(c *Container) (interface{}, error) {
	return c.Get(a.ref)
}

func ServiceArg(ref fmt.Stringer) ServiceDefArg {
	return &serviceRefArg{ref: ref}
}

// Service Method Call Argument allows to use the return value of a method call on given service.
type serviceMethodCallArg struct {
	serviceRef fmt.Stringer
	methodName string
	args       []ServiceDefArg
}

func (a *serviceMethodCallArg) evaluate(c *Container) (interface{}, error) {
	s, err := c.Get(a.serviceRef)
	if err != nil {
		return nil, err
	}

	return c.callReflectValueWithArgs(reflect.ValueOf(s).MethodByName(a.methodName), a.args)
}

func ServiceMethodCallArg(serviceRef fmt.Stringer, methodName string, args ...ServiceDefArg) ServiceDefArg {
	return &serviceMethodCallArg{
		serviceRef: serviceRef,
		methodName: methodName,
		args:       args,
	}
}

// Services By Tag Arg allows to get one or more services by tags and inject them.
type servicesByTagArg struct {
	tags []fmt.Stringer
}

func (a *servicesByTagArg) evaluate(c *Container) (interface{}, error) {
	return c.FindByTags(a.tags)
}

// ServicesByTagsArg is a shortcut for a service argument.
//goland:noinspection GoUnusedExportedFunction
func ServicesByTagsArg(tags []fmt.Stringer) ServiceDefArg {
	return &servicesByTagArg{tags: tags}
}

// Parameter Argument allows to get parameters by path/dot notation from parameter provider.
type paramArg struct {
	paramPath string
}

func (a *paramArg) evaluate(c *Container) (interface{}, error) {
	return c.paramProvider.Get(a.paramPath)
}

func ParamArg(paramPath string) ServiceDefArg {
	return &paramArg{paramPath: paramPath}
}

// Context Argument injects the context from di.Container.
type contextArg struct{}

func (a *contextArg) evaluate(c *Container) (interface{}, error) {
	return c.ctx, nil
}

func ContextArg() ServiceDefArg {
	return &contextArg{}
}

// Container Argument injects the container from di.Container.
type containerArg struct{}

func (a *containerArg) evaluate(c *Container) (interface{}, error) {
	return c, nil
}

func ContainerArg() ServiceDefArg {
	return &containerArg{}
}

// EventBus Argument injects the eventBus from di.EventBus.
type eventBusArg struct{}

func (a *eventBusArg) evaluate(c *Container) (interface{}, error) {
	return c.eventBus, nil
}

//goland:noinspection GoUnusedExportedFunction
func EventBusArg() ServiceDefArg {
	return &eventBusArg{}
}
