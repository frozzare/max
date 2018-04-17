package cast

import (
	"errors"
	"strconv"
)

// Bool converts argument to bool or return a error.
func Bool(value interface{}) (bool, error) {
	switch v := value.(type) {
	case bool:
		return v, nil
	case float32:
		return v > 0, nil
	case float64:
		return v > 0, nil
	case int:
		return v > 0, nil
	case int8:
		return v > 0, nil
	case int16:
		return v > 0, nil
	case int32:
		return v > 0, nil
	case int64:
		return v > 0, nil
	case uint:
		return v > 0, nil
	case uint8:
		return v > 0, nil
	case uint16:
		return v > 0, nil
	case uint32:
		return v > 0, nil
	case uint64:
		return v > 0, nil
	case string:
		return strconv.ParseBool(v)
	case []byte:
		return strconv.ParseBool(string(v))
	case nil:
		return false, nil
	default:
		return false, errors.New("Unknown type")
	}
}

// MustBool converts argument to bool or panic if an error occurred.
func MustBool(value interface{}) bool {
	v, err := Bool(value)

	if err != nil {
		panic(err)
	}

	return v
}
