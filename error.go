package di

//go:generate stringer -type=ErrType

type ErrType int

const (
	// Simple error codes.
	ErrServiceNotFound ErrType = iota
	ErrProviderMissing
	ErrProviderNotAFunc
	ErrProviderToManyReturnValues
	ErrProviderArgCountMismatch
	ErrProviderArgTypeMismatch
	ErrParamProviderGet
	ErrParamProviderNotDefined
)
