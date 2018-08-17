package template2

import (
	"reflect"
)

// Isset determine if a map key or struct field is set and valid.
func Isset(data interface{}, name string) bool {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() == reflect.Map {
		for _, key := range v.MapKeys() {
			if key.String() == name && key.IsValid() {
				return v.MapIndex(key).IsValid()
			}
		}
	}

	if v.Kind() != reflect.Struct {
		return false
	}

	return v.FieldByName(name).IsValid()
}
