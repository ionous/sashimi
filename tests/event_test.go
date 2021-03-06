package tests

import (
	. "github.com/ionous/sashimi/script"
	"github.com/ionous/sashimi/standard"
	//	"github.com/ionous/sashimi/standard/framework"
	G "github.com/ionous/sashimi/game"
	"github.com/stretchr/testify/assert"
	"testing"
)

// FIX: not really unit test yet, but helpful for debugging
// really, this should have an event sink mock and test that everything works as expected.
func TestEvent(t *testing.T) {
	visited := false
	s := EventScript(t, &visited)
	if test, err := NewTestGame(t, s); assert.NoError(t, err) {

		if here, ok := test.Game.GetInstance("here"); assert.True(t, ok) {
			here := test.Game.NewAdapter().NewGameObject(here)
			assert.False(t, visited)
			here.Go("report the view")
			if r, e := test.FlushOutput(); assert.NoError(t, e) {
				t.Log(r)
				assert.True(t, visited)
			}
		}
	}
}

func EventScript(t *testing.T, visited *bool) *Script {
	s := standard.InitStandardLibrary()
	s.The("room",
		Called("here"),
		Has("description", "a dull room"),
		After("reporting the view").Always(func(g G.Play) {
			t.Log("after reporting the view")
			*visited = true
		}),
	)
	s.The("object", Called("the knicknack"),
		Exists(),
		In("here"))
	s.The("player",
		Exists(),
		In("here"))
	return s
}
