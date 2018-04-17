package slices

import (
	"testing"

	"github.com/frozzare/go-assert"
)

func TestMap(t *testing.T) {
	v, err := Map([]int{1, 2, 3}, func(v int) int {
		return v + 1
	})

	assert.Nil(t, err)
	assert.Equal(t, []int{2, 3, 4}, v.([]int))

	v2, err := Map([]int{1, 2, 3}, func(v, i int) int {
		return v + i
	})

	assert.Nil(t, err)
	assert.Equal(t, []int{1, 3, 5}, v2.([]int))
}
