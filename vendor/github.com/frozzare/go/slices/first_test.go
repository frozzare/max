package slices

import (
	"testing"

	"github.com/frozzare/go-assert"
)

func TestFirst(t *testing.T) {
	slice := []string{"fredrik", "elli", "go"}

	assert.Equal(t, "fredrik", First(slice))

	assert.Nil(t, First([]string{}))
	assert.Nil(t, First(1))
}
