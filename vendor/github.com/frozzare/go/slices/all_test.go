package slices

import (
	"testing"

	"github.com/frozzare/go-assert"
)

func TestAll(t *testing.T) {
	v := All([]int{2, 3, 1}, func(v int) bool {
		return v > 0
	})

	assert.True(t, v)
}
