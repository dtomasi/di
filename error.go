package di

//go:generate stringer -type=ErrType

type ErrType int

const (
	// Simple error codes.
	ErrServiceNotFound ErrType = iota
	ErrProviderMissing
	ErrProviderNotAFunc
	ErrParamProviderGet
	ErrProviderToManyReturnValues
	ErrProviderArgCountMismatch
	ErrProviderArgTypeMismatch
)
