package structs

import (
	"errors"
	"fmt"
	"reflect"
)

const defaultTag = "structs"

func strct(s interface{}) (reflect.Value, error) {
	v := reflect.ValueOf(s)

	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return v, fmt.Errorf("%s is not a struct", v.Kind())
	}

	return v, nil
}

func fields(s interface{}, args ...string) ([]*StructField, error) {
	name := defaultTag
	if len(args) > 0 {
		name = args[0]
	}

	v, err := strct(s)

	if err != nil {
		return make([]*StructField, 0), err
	}

	t := v.Type()

	fields := make([]*StructField, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		if tag := f.Tag.Get(name); tag == "-" {
			continue
		}

		fields[i] = &StructField{
			name:  f.Name,
			value: v.FieldByName(f.Name),
		}
	}

	return fields, nil
}

// Field returns a field from a struct by name.
func Field(s interface{}, n string) (*StructField, error) {
	r, err := strct(s)
	if err != nil {
		return nil, err
	}

	field, ok := r.Type().FieldByName(n)
	if !ok {
		return nil, errors.New("Field not found")
	}

	return &StructField{
		field: field,
		name:  n,
		value: r.FieldByName(n),
	}, nil
}

// Fields returns a list field. Tag argument is optional.
func Fields(s interface{}, tag ...string) ([]*StructField, error) {
	return fields(s, tag...)
}

// Name returns the structs type name within it's package.
// Empty string if not a struct.
func Name(s interface{}) (string, error) {
	r, err := strct(s)

	if err != nil {
		return "", err
	}

	return r.Type().Name(), err
}

// Names returns the struct fields name.
// Field is ignored with tag `structs:"-"` or tag of your choice.
// Empty slice if not a struct.
func Names(s interface{}, tag ...string) ([]string, error) {
	fields, err := fields(s, tag...)

	if err != nil {
		return make([]string, 0), err
	}

	names := make([]string, len(fields))

	for i := 0; i < len(fields); i++ {
		names[i] = fields[i].Name()
	}

	return names, nil
}
