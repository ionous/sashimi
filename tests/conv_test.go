package tests

import (
	_ "github.com/ionous/sashimi/extensions"
	. "github.com/ionous/sashimi/script"
	"github.com/stretchr/testify/assert"
	"testing"
)

//
func TestConversationCompilation(t *testing.T) {
	s := InitScripts()
	_, err := NewTestGame(t, s)
	assert.NoError(t, err)
}
