// Code generated by "stringer -type=ArgType"; DO NOT EDIT.

package di

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ArgTypeInterface-0]
	_ = x[ArgTypeContext-1]
	_ = x[ArgTypeContainer-2]
	_ = x[ArgTypeService-3]
	_ = x[ArgTypeServicesByTags-4]
	_ = x[ArgTypeParam-5]
}

const _ArgType_name = "ArgTypeInterfaceArgTypeContextArgTypeContainerArgTypeServiceArgTypeServicesByTagsArgTypeParam"

var _ArgType_index = [...]uint8{0, 16, 30, 46, 60, 81, 93}

func (i ArgType) String() string {
	if i < 0 || i >= ArgType(len(_ArgType_index)-1) {
		return "ArgType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ArgType_name[_ArgType_index[i]:_ArgType_index[i+1]]
}
