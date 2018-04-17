package cast

import (
	"testing"

	"github.com/frozzare/go-assert"
)

func TestFloat32(t *testing.T) {
	test := [][]interface{}{
		{'f', 102.0},
		{true, 1.0},
		{false, 0.0},
		{int(0), 0.0},
		{int8(1), 1.0},
		{int16(16), 16.0},
		{int32(32), 32.0},
		{int64(64), 64.0},
		{uint(0), 0.0},
		{uint8(1), 1.0},
		{uint16(16), 16.0},
		{uint32(32), 32.0},
		{uint64(64), 64.0},
		{float32(0.1), 0.10000000149011612},
		{float64(1.1), 1.1},
		{float64(1.348959), 1.348959},
		{"1.2", 1.2},
		{[]byte("1.2"), 1.2000000476837158},
		// errors
		{"x", 0.0, true},
	}

	for _, item := range test {
		v, err := Float32(item[0])

		if len(item) > 2 {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}

		assert.Equal(t, float32(item[1].(float64)), v)
	}
}

func TestFloat64(t *testing.T) {
	test := [][]interface{}{
		{'f', 102.0},
		{true, 1.0},
		{false, 0.0},
		{int(0), 0.0},
		{int8(1), 1.0},
		{int16(16), 16.0},
		{int32(32), 32.0},
		{int64(64), 64.0},
		{uint(0), 0.0},
		{uint8(1), 1.0},
		{uint16(16), 16.0},
		{uint32(32), 32.0},
		{uint64(64), 64.0},
		{float32(0.1), 0.10000000149011612},
		{float64(1.1), 1.1},
		{float64(1.348959), 1.348959},
		{"1.2", 1.2},
		{[]byte("1.2"), 1.2000000476837158},
		// errors
		{"x", 0.0, true},
	}

	for _, item := range test {
		v, err := Float64(item[0])

		if len(item) > 2 {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}

		assert.Equal(t, item[1].(float64), v)
	}
}
