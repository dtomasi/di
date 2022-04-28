// Code generated by "stringer -type=ErrorType"; DO NOT EDIT.

package di

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ContainerBuildError-0]
	_ = x[ServiceNotFoundError-1]
	_ = x[ServiceBuildError-2]
	_ = x[ProviderMissingError-3]
	_ = x[ProviderNotAFuncError-4]
	_ = x[ProviderToManyReturnValuesError-5]
	_ = x[ProviderArgCountMismatchError-6]
	_ = x[ProviderArgTypeMismatchError-7]
	_ = x[ParamProviderNotDefinedError-8]
}

const _ErrorType_name = "ContainerBuildErrorServiceNotFoundErrorServiceBuildErrorProviderMissingErrorProviderNotAFuncErrorProviderToManyReturnValuesErrorProviderArgCountMismatchErrorProviderArgTypeMismatchErrorParamProviderNotDefinedError"

var _ErrorType_index = [...]uint8{0, 19, 39, 56, 76, 97, 128, 157, 185, 213}

func (i ErrorType) String() string {
	if i < 0 || i >= ErrorType(len(_ErrorType_index)-1) {
		return "ErrorType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ErrorType_name[_ErrorType_index[i]:_ErrorType_index[i+1]]
}
