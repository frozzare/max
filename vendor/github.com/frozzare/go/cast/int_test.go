package cast

import (
	"testing"

	"github.com/frozzare/go-assert"
)

func TestInt(t *testing.T) {
	test := [][]interface{}{
		{'f', 102},
		{true, 1},
		{false, 0},
		{int(0), 0},
		{int8(1), 1},
		{int16(16), 16},
		{int32(32), 32},
		{int64(64), 64},
		{uint(0), 0},
		{uint8(1), 1},
		{uint16(16), 16},
		{uint32(32), 32},
		{uint64(64), 64},
		{float32(0.1), 0},
		{float64(1.1), 1},
		{float64(1.348959), 1},
		{"1", 1},
		{"1.2", 1},
		{[]byte("1.2"), 1},
		{"x", 0},
	}

	for _, item := range test {
		v, err := Int(item[0])
		assert.Nil(t, err)
		assert.Equal(t, item[1].(int), v)
	}
}

func TestInt64(t *testing.T) {
	test := [][]interface{}{
		{'f', 102},
		{true, 1},
		{false, 0},
		{int(0), 0},
		{int8(1), 1},
		{int16(16), 16},
		{int32(32), 32},
		{int64(64), 64},
		{uint(0), 0},
		{uint8(1), 1},
		{uint16(16), 16},
		{uint32(32), 32},
		{uint64(64), 64},
		{float32(0.1), 0},
		{float64(1.1), 1},
		{float64(1.348959), 1},
		{"1", 1},
		{"1.2", 1},
		{[]byte("1.2"), 1},
		{"x", 0},
	}

	for _, item := range test {
		v, err := Int64(item[0])
		assert.Nil(t, err)
		assert.Equal(t, item[1].(int), v)
	}
}
