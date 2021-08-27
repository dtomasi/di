package di

import "fmt"

// ArgType defines the type of argument.
// Currently, those are service or param (parameter from service container).
type ArgType int

const (
	ArgTypeInterface ArgType = iota
	ArgTypeContext
	ArgTypeContainer
	ArgTypeParamProvider
	ArgTypeService
	ArgTypeParam
)

// Arg defines a argument for provider functions or defined calls.
type Arg struct {
	_type ArgType
	value interface{}
}

// ArgWithType defines an argument by type, name and value.
func ArgWithType(argType ArgType, argValue interface{}) Arg {
	return Arg{
		_type: argType,
		value: argValue,
	}
}

// InterfaceArg is a shortcut for an argument of type interface{}.
func InterfaceArg(in interface{}) Arg {
	return ArgWithType(ArgTypeInterface, in)
}

// ServiceArg is a shortcut for a service argument.
func ServiceArg(serviceRef fmt.Stringer) Arg {
	return ArgWithType(ArgTypeService, serviceRef)
}

// ParamArg is a shortcut for a parameter argument.
func ParamArg(paramPath string) Arg {
	return ArgWithType(ArgTypeParam, paramPath)
}

// ContextArg is a shortcut for an argument with no value that injects the context.
func ContextArg() Arg {
	return ArgWithType(ArgTypeContext, nil)
}

// ContainerArg is a shortcut for an argument with no value that injects the container itself.
func ContainerArg() Arg {
	return ArgWithType(ArgTypeContainer, nil)
}

// ParamProviderArg is a shortcut for an argument with no value that injects the container´s parameter provider.
func ParamProviderArg() Arg {
	return ArgWithType(ArgTypeParamProvider, nil)
}
