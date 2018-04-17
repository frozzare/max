package slices

import (
	"testing"

	"github.com/frozzare/go-assert"
)

func TestReject(t *testing.T) {
	v, err := Reject([]int{1, 2, 3}, func(v int) bool {
		return v%2 == 0
	})

	assert.Nil(t, err)
	assert.Equal(t, []int{1, 3}, v.([]int))
}
