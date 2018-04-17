package slices

import (
	"testing"

	"github.com/frozzare/go-assert"
)

func TestLast(t *testing.T) {
	slice := []string{"fredrik", "elli", "go"}

	assert.Equal(t, "go", Last(slice))

	assert.Nil(t, Last([]string{}))
	assert.Nil(t, Last(1))
}
