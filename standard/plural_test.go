package standard

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// go test --run TestVowels
func TestVowels(t *testing.T) {
	assert.True(t, startsVowel("evil fish"))
}
