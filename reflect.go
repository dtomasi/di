package di

import "reflect"

// getType is a simple function for getting type as string for comparison.
func getType(ty reflect.Type) string {
	t := ty
	if t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	}

	return t.Name()
}
