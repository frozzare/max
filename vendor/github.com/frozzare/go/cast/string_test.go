package cast

import (
	"testing"

	"github.com/frozzare/go-assert"
)

func TestString(t *testing.T) {
	test := [][]interface{}{
		{[]byte("f"), "f"},
		{true, "true"},
		{false, "false"},
		{int(0), "0"},
		{int8(1), "1"},
		{int16(16), "16"},
		{int32(32), "32"},
		{int64(64), "64"},
		{uint(0), "0"},
		{uint8(1), "1"},
		{uint16(16), "16"},
		{uint32(32), "32"},
		{uint64(64), "64"},
		{float32(0.1), "0.10000000149011612"},
		{float64(1.1), "1.1"},
		{float64(1.348959), "1.348959"},
		{"hello", "hello"},
		{[]int{1, 2, 3}, "[1 2 3]"},
		{nil, ""},
		{[]byte("hello"), "hello"},
	}

	for _, item := range test {
		v, err := String(item[0])
		assert.Nil(t, err)
		assert.Equal(t, item[1].(string), v)
	}
}
