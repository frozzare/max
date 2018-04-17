package structs

import (
	"errors"
	"reflect"

	"github.com/frozzare/go/reflect2"
)

// StructField represents a field.
type StructField struct {
	field reflect.StructField
	name  string
	value reflect.Value
}

// IsZero reports whether a value is a zero value of its kind.
func (s *StructField) IsZero() bool {
	return reflect2.IsZero(s.value)
}

// Tag returns a field tag value by key or a empty string.
func (s *StructField) Tag(k string) string {
	return s.field.Tag.Get(k)
}

// Kind returns the fields kind.
func (s *StructField) Kind() reflect.Kind {
	return s.value.Kind()
}

// Name returns the field name.
func (s *StructField) Name() string {
	return s.name
}

// Set sets the field value only if field is exported and can be set.
func (s *StructField) Set(v interface{}) error {
	if s.field.PkgPath != "" {
		return errors.New("Field is not exported")
	}

	if !s.value.CanSet() {
		return errors.New("Field cannot be set")
	}

	given := reflect.ValueOf(v)

	if s.value.Kind() != given.Kind() {
		return errors.New("Field and value kind don't match")
	}

	s.value.Set(given)
	return nil
}

// Value returns the field value.
func (s *StructField) Value() interface{} {
	return s.value.Interface()
}
