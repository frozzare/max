package slices

import (
	"testing"

	"github.com/frozzare/go-assert"
)

func TestFilter(t *testing.T) {
	v, err := Filter([]int{1, 2, 3}, func(v int) bool {
		return v%2 == 0
	})

	assert.Nil(t, err)
	assert.Equal(t, []int{2}, v.([]int))
}
