package cast

import (
	"testing"

	"github.com/frozzare/go-assert"
)

func TestBool(t *testing.T) {
	test := [][]interface{}{
		{'f', true},
		{true, true},
		{false, false},
		{int(0), false},
		{int8(1), true},
		{int16(16), true},
		{int32(32), true},
		{int64(64), true},
		{uint(0), false},
		{uint8(1), true},
		{uint16(16), true},
		{uint32(32), true},
		{uint64(64), true},
		{float32(0.1), true},
		{float64(1.1), true},
		{float64(1.348959), true},
		{"true", true},
		{"false", false},
		{"1", true},
		{"0", false},
		{[]byte("true"), true},
		{[]byte("false"), false},
		// errors
		{"x", false, true},
		{"1.5", false, true},
	}

	for _, item := range test {
		v, err := Bool(item[0])

		if len(item) > 2 {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}

		assert.Equal(t, item[1].(bool), v)
	}
}
