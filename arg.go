package di

import "fmt"

type ServiceDefArg interface {
	Evaluate(*Container) (interface{}, error)
}

// Interface Arg
// This argument allows to pass any value that is provided to the service.
type interfaceArg struct {
	inValue interface{}
}

func (a *interfaceArg) Evaluate(_ *Container) (interface{}, error) {
	return a.inValue, nil
}

// InterfaceArg is a shortcut for an argument of type interface{}.
func InterfaceArg(in interface{}) ServiceDefArg {
	return &interfaceArg{inValue: in}
}

// ServiceRef Arg
// This argument type allows to pass a reference to a service that will be injected.
type serviceRefArg struct {
	ref fmt.Stringer
}

func (a *serviceRefArg) Evaluate(c *Container) (interface{}, error) {
	return c.Get(a.ref)
}

func ServiceRefArg(ref fmt.Stringer) ServiceDefArg {
	return &serviceRefArg{ref: ref}
}

// Deprecated: use ServiceRefArg.
func ServiceArg(ref fmt.Stringer) ServiceDefArg {
	return ServiceRefArg(ref)
}

// Services By Tag Arg allows to get one or more services by tags and inject them.
type servicesByTagArg struct {
	tags []fmt.Stringer
}

func (a *servicesByTagArg) Evaluate(c *Container) (interface{}, error) {
	return c.FindByTags(a.tags)
}

// ServicesByTagsArg is a shortcut for a service argument.
func ServicesByTagsArg(tags []fmt.Stringer) ServiceDefArg {
	return &servicesByTagArg{tags: tags}
}

// Parameter Argument allows to get parameters by path/dot notation from parameter provider.
type paramArg struct {
	paramPath string
}

func (a *paramArg) Evaluate(c *Container) (interface{}, error) {
	return c.paramProvider.Get(a.paramPath)
}

func ParamArg(paramPath string) ServiceDefArg {
	return &paramArg{paramPath: paramPath}
}

// Context Argument injects the context from di.Container.
type contextArg struct{}

func (a *contextArg) Evaluate(c *Container) (interface{}, error) {
	return c.ctx, nil
}

func ContextArg() ServiceDefArg {
	return &contextArg{}
}

// Container Argument injects the container from di.Container.
type containerArg struct{}

func (a *containerArg) Evaluate(c *Container) (interface{}, error) {
	return c, nil
}

func ContainerArg() ServiceDefArg {
	return &containerArg{}
}

// EventBus Argument injects the eventBus from di.EventBus.
type eventBusArg struct{}

func (a *eventBusArg) Evaluate(c *Container) (interface{}, error) {
	return c.eventBus, nil
}

func EventBusArg() ServiceDefArg {
	return &eventBusArg{}
}
