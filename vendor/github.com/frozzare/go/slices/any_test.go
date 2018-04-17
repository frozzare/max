package slices

import (
	"testing"

	"github.com/frozzare/go-assert"
)

func TestAny(t *testing.T) {
	v := Any([]int{2, 3, 1}, func(v int) bool {
		return v == 1
	})

	assert.True(t, v)
}
