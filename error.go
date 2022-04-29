package di

//go:generate stringer -type=ErrorType -output=zz_gen_errortype_string.go

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
