package di

import (
	"fmt"
	"reflect"
)

// StringRef defines a service reference as string.
// As ServiceDef implements fmt.Stringer as an interface for referencing services
type StringRef string

// String implements fmt.Stringer interface method
func (r StringRef) String() string {
	return string(r)
}

// ServiceDefOption defines an option function
type ServiceDefOption func(def *ServiceDef)

// BuildOnFirstRequest option will not create an instance from on building the container
// Instead it will wait until the service is requested the first time
func BuildOnFirstRequest() ServiceDefOption {
	return func(sd *ServiceDef) {
		sd.options.buildOnFirstRequest = true
	}
}

// serviceDefOptions lifecycle options for definition
type serviceDefOptions struct {
	buildOnFirstRequest bool
}

// ServiceDef is a definition of a service
// it describes how a service should be created and handled inside the
// service container
type ServiceDef struct {
	ref      fmt.Stringer
	options  serviceDefOptions
	provider interface{}
	args     []Arg
}

// NewServiceDef creates a new service definition
func NewServiceDef(ref fmt.Stringer) *ServiceDef {
	i := &ServiceDef{
		ref: ref,
		options: serviceDefOptions{
			buildOnFirstRequest: false,
		},
	}
	return i
}

// Opts allows to set some options for lifecycle management
func (sd *ServiceDef) Opts(opts... ServiceDefOption) *ServiceDef {
	for _, opt := range opts {
		opt(sd)
	}
	return sd
}

// Provider defines a function that returns the actual serve instance
// this function can also accept arguments that are described using the Args function
func (sd *ServiceDef) Provider(provider interface{}) *ServiceDef {
	sd.provider = provider
	return sd
}

// Args accepts multiple constructor/provider function arguments
func (sd *ServiceDef) Args(args ...Arg) *ServiceDef {
	sd.args = append(sd.args, args...)
	return sd
}

// build creates the service instance using passed parsed arguments from container
func (sd *ServiceDef) build(parsedArgs []Arg) (interface{}, error) {
	sd.args = parsedArgs

	// build using provider function
	if sd.provider != nil {
		x := reflect.TypeOf(sd.provider)

		if x.Kind() != reflect.Func {
			return nil, fmt.Errorf("provider is not a func type")
		}

		inputArgCount := x.NumIn()
		if inputArgCount != len(sd.args) {
			return nil, fmt.Errorf(
				"expected %d arguments for %s provider. Got %d",
				inputArgCount,
				sd.ref,
				len(sd.args),
			)
		}

		var inputValues []reflect.Value
		for i := 0; i < inputArgCount; i++ {
			inType := x.In(i)

			inArgType := reflect.TypeOf(sd.args[i].value)
			if getType(inType) != getType(inArgType) && !inArgType.Implements(inType) {
				return nil, fmt.Errorf(
					"provider argument at position %d should be type of or implementing %s. Got %s",
					i,
					getType(inType),
					getType(inArgType),
				)
			}

			inputValues = append(inputValues, reflect.ValueOf(sd.args[i].value))
		}

		y := reflect.ValueOf(sd.provider)
		returnValues := y.Call(inputValues)

		switch len(returnValues) {
		case 1:
			return returnValues[0].Interface(), nil
		case 2:
			err, ok := returnValues[1].Interface().(error)
			if !ok {
				err = nil
			}
			return returnValues[0].Interface(), err
		default:
			return nil, fmt.Errorf("to many return values in provider function for services %s", sd.ref)
		}
	}

	return nil, fmt.Errorf("no provider function set for service %s", sd.ref)
}

// ArgType defines the type of argument.
// Currently, those are service or param (parameter from service container)
type ArgType int

const (
	ArgTypeInterface ArgType = iota
	ArgTypeContext
	ArgTypeContainer
	ArgTypeParamProvider
	ArgTypeService
	ArgTypeParam
)

// Arg defines a argument for provider functions or defined calls
type Arg struct {
	_type ArgType
	value interface{}
}

// ArgWithType defines an argument by type, name and value
func ArgWithType(argType ArgType, argValue interface{}) Arg {
	return Arg{
		_type: argType,
		value: argValue,
	}
}

// InterfaceArg is a shortcut for an argument of type interface{}
func InterfaceArg(in interface{}) Arg {
	return ArgWithType(ArgTypeInterface, in)
}

// ServiceArg is a shortcut for a service argument
func ServiceArg(serviceRef fmt.Stringer) Arg {
	return ArgWithType(ArgTypeService, serviceRef)
}

// ParamArg is a shortcut for a parameter argument
func ParamArg(paramPath string) Arg {
	return ArgWithType(ArgTypeParam, paramPath)
}

// ContextArg is a shortcut for an argument with no value that injects the context
func ContextArg() Arg {
	return ArgWithType(ArgTypeContext, nil)
}

// ContainerArg is a shortcut for an argument with no value that injects the container itself
func ContainerArg() Arg {
	return ArgWithType(ArgTypeContainer, nil)
}

// ParamProviderArg is a shortcut for an argument with no value that injects the containers parameter provider
func ParamProviderArg() Arg {
	return ArgWithType(ArgTypeParamProvider, nil)
}

// getType is a simple function for getting type as string for comparison
func getType(ty reflect.Type) string {
	if t := ty; t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}
