package slices

import (
	"reflect"
)

func source(source interface{}) (int, func(int) (reflect.Value, reflect.Value)) {
	if source == nil {
		return 0, nil
	}

	v := reflect.ValueOf(source)
	switch v.Kind() {
	case reflect.Array:
	case reflect.Slice:
		return v.Len(), func(i int) (reflect.Value, reflect.Value) {
			return reflect.ValueOf(i), v.Index(i)
		}
	case reflect.Map:
		mk := v.MapKeys()
		return len(mk), func(i int) (reflect.Value, reflect.Value) {
			return mk[i], v.MapIndex(mk[i])
		}
	}

	return 0, nil
}

func each(input, iterator interface{}, predicate func(reflect.Value, reflect.Value, reflect.Value) bool) {
	l, item := source(input)

	if l == 0 {
		return
	}

	if predicate == nil {
		predicate = func(output, _, _ reflect.Value) bool {
			if output.Kind() == reflect.Bool {
				return output.Bool()
			}

			return false
		}
	}

	iv := reflect.ValueOf(iterator)
	it := reflect.TypeOf(iterator)
	lin := it.NumIn()

	for i := 0; i < l; i++ {
		k, v := item(i)

		in := []reflect.Value{}

		if lin == 2 {
			in = append(in, v)
			in = append(in, k)
		} else {
			in = append(in, v)
		}

		current := iv.Call(in)

		if len(current) > 0 && predicate(current[0], k, v) {
			break
		}
	}
}
