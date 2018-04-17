package slices

import (
	"testing"

	"github.com/frozzare/go-assert"
)

func TestUniq(t *testing.T) {
	v, err := Uniq([]int{1, 2, 3, 1, 2, 3}, nil)

	assert.Nil(t, err)
	assert.Equal(t, []int{1, 2, 3}, v.([]int))

	v2, err := Uniq([]string{"fredrik", "elli", "farfar", "elli"}, func(s string) bool {
		return s[0] == 'f'
	})

	assert.Nil(t, err)
	assert.Equal(t, []string{"fredrik", "farfar"}, v2.([]string))
}
