package di

//go:generate stringer -type=ErrorType

type ErrorType int

const (
	// Simple error codes.
	ContainerBuildError ErrorType = iota
	ServiceNotFoundError
	ServiceBuildError
	ProviderMissingError
	CallableNotAFuncError
	CallableToManyReturnValuesError
	CallableArgCountMismatchError
	CallableArgTypeMismatchError
	ParamProviderNotDefinedError
)
