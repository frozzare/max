package cast

import (
	"errors"
	"strconv"
)

// Float32 will converts argument to float32 or return a error.
func Float32(value interface{}) (float32, error) {
	switch v := value.(type) {
	case bool:
		if v {
			return float32(1), nil
		}

		return float32(0), nil
	case float32:
		return float32(v), nil
	case float64:
		return float32(v), nil
	case int:
		return float32(v), nil
	case int8:
		return float32(v), nil
	case int16:
		return float32(v), nil
	case int32:
		return float32(v), nil
	case int64:
		return float32(v), nil
	case uint:
		return float32(v), nil
	case uint8:
		return float32(v), nil
	case uint16:
		return float32(v), nil
	case uint32:
		return float32(v), nil
	case uint64:
		return float32(v), nil
	case string:
		f, err := strconv.ParseFloat(v, 64)
		return float32(f), err
	case []byte:
		f, err := strconv.ParseFloat(string(v), 32)
		return float32(f), err
	case nil:
		return float32(0), nil
	default:
		return float32(0), errors.New("Unknown type")
	}
}

// MustFloat32 converts argument to float32 or panic if an error occurred.
func MustFloat32(value interface{}) float32 {
	v, err := Float32(value)

	if err != nil {
		panic(err)
	}

	return v
}

// Float64 will converts argument to float64 or return a error.
func Float64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case bool:
		if v {
			return float64(1), nil
		}

		return float64(0), nil
	case float32:
		return float64(v), nil
	case float64:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case string:
		f, err := strconv.ParseFloat(v, 64)
		return float64(f), err
	case []byte:
		f, err := strconv.ParseFloat(string(v), 32)
		return float64(f), err
	case nil:
		return float64(0), nil
	default:
		return float64(0), errors.New("Unknown type")
	}
}

// MustFloat64 converts argument to float64 or panic if an error occurred.
func MustFloat64(value interface{}) float64 {
	v, err := Float64(value)

	if err != nil {
		panic(err)
	}

	return v
}
