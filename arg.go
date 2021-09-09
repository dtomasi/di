package di

import "fmt"

//go:generate stringer -type=ArgType

// ArgType defines the type of argument.
// Currently, those are service or param (parameter from service container).
type ArgType int

const (
	ArgTypeInterface ArgType = iota
	ArgTypeContext
	ArgTypeContainer
	ArgTypeParamProvider
	ArgTypeLogger
	ArgTypeNamedLogger
	ArgTypeService
	ArgTypeParam
)

type ArgParseEvent struct {
	ServiceRef fmt.Stringer
	Pos int
	Arg Arg
	Err error
}

// Arg defines a argument for provider functions or defined calls.
type Arg struct {
	_type ArgType
	value interface{}
}

func (a Arg) GetType() ArgType {
	return a._type
}

func (a Arg) GetValue() interface{} {
	return a.value
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

// LoggerArg is a shortcut for an argument with no value that injects the logger provided with the container.
func LoggerArg() Arg {
	return ArgWithType(ArgTypeLogger, nil)
}

// LoggerArg is a shortcut for an argument with no value that injects the logger provided with the container.
func NamedLoggerArg(name string) Arg {
	return ArgWithType(ArgTypeNamedLogger, name)
}

// ParamProviderArg is a shortcut for an argument with no value that injects the container´s parameter provider.
func ParamProviderArg() Arg {
	return ArgWithType(ArgTypeParamProvider, nil)
}
