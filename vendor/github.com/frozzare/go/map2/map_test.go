package map2

import (
	"testing"

	"github.com/frozzare/go-assert"
)

var m = map[string]int{
	"a": 1,
	"b": 2,
	"c": 3,
	"d": 4,
	"e": 5,
	"f": 6,
}

var m2 = map[interface{}]interface{}{
	1:   "2",
	"4": 3,
}

func TestKeys(t *testing.T) {
	v, err := Keys(m)
	assert.Nil(t, err)
	assert.Equal(t, 6, len(v.([]string)))

	v, err = Keys(m2)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(v.([]interface{})))
}

func TestValues(t *testing.T) {
	v, err := Values(m)
	assert.Nil(t, err)
	assert.Equal(t, 6, len(v.([]int)))

	v, err = Keys(m2)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(v.([]interface{})))
}
