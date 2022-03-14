package utils

import "reflect"

// GetType is a simple function for getting type as string for comparison.
func GetType(ty reflect.Type) string {
	t := ty
	if t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	}

	return t.Name()
}
