package slices

import (
	"testing"

	"github.com/frozzare/go-assert"
)

func TestFold(t *testing.T) {
	slice := []int{1, 2, 3}

	v, _ := Fold(slice, func(v1, v2 int) int {
		return v1 + v2
	})

	assert.Equal(t, 6, v)

	slice2 := []string{"1", "2", "3"}

	v2, _ := Fold(slice2, func(v1, v2 string) string {
		return v1 + v2
	})

	assert.Equal(t, "123", v2)
}
