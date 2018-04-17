package net2

import (
	"testing"

	"github.com/frozzare/go-assert"
)

func TestRandomPort(t *testing.T) {
	p, err := RandomPort()
	assert.Nil(t, err)
	assert.Equal(t, true, p > 0)
}
