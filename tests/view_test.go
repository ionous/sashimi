package tests

import (
	. "github.com/ionous/sashimi/script"
	"github.com/ionous/sashimi/standard/framework"
	. "github.com/ionous/sashimi/standard/live"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestViewRoot(t *testing.T) {
	s := ViewScript()
	if test, err := NewTestGame(t, s); assert.NoError(t, err) {
		if here, ok := test.Game.GetInstance("here"); assert.True(t, ok) {
			if there, ok := test.Game.GetInstance("there"); assert.True(t, ok) {
				if player, ok := test.Game.GetInstance("player"); assert.True(t, ok) {
					p := framework.NewParentLookup(test.Game)
					assert.Equal(t, here, p.LookupRoot(here), "root is here")
					assert.Equal(t, here, p.LookupRoot(player), "player is here")
					assert.Equal(t, there, p.LookupRoot(there), "root is there")
				}
			}
		}
	}
}
func TestViewChange(t *testing.T) {
	s := ViewScript()
	if test, err := NewTestGame(t, s); assert.NoError(t, err) {
		if here, ok := test.Game.GetInstance("here"); assert.True(t, ok) {
			if there, ok := test.Game.GetInstance("there"); assert.True(t, ok) {
				if player, ok := test.Game.GetInstance("player"); assert.True(t, ok) {
					view := framework.NewStandardView(test.Game)
					//
					assert.Equal(t, player.GetId(), view.Viewer(), "player is viewpoint")
					assert.True(t, view.InView(player), "viewpoint in view")
					//
					assert.True(t, view.InView(here), "here in view")
					assert.False(t, view.InView(there), "there not in view")
					//
					assert.True(t, view.ChangedView(player, "objects-whereabouts", there), "view changed")
					//
					assert.True(t, view.InView(there), "there in view")
					assert.False(t, view.InView(here), "here no longer in view")
					// move into there
					test.Game.NewAdapter().Go(Move("the player").To("there"))
					assert.NoError(t, test.Game.ProcessActions(), "moved")
					assert.Equal(t, player.GetId(), view.Viewer(), "player is still viewpoint")
					assert.True(t, view.InView(player), "viewpoint still in view")
				}
			}
		}
	}
}

func ViewScript() *Script {
	s := InitScripts()
	s.The("room",
		Called("here"),
		Has("description", "an empty room"),
	)
	s.The("room",
		Called("there"),
		Has("description", "an empty room"),
	)
	s.The("player",
		Exists(),
		In("here"))
	return s
}
