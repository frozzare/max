package strings2

import (
	"testing"

	"github.com/frozzare/go-assert"
)

func TestSwapCase(t *testing.T) {
	assert.Equal(t, "FredriK", SwapCase("fREDRIk"))
	assert.Equal(t, "Åääö", SwapCase("åÄÄÖ"))
	assert.Equal(t, "hELLO, 世界", SwapCase("Hello, 世界"))
}
