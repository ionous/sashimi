package tests

import (
	G "github.com/ionous/sashimi/game"
	R "github.com/ionous/sashimi/runtime"
	. "github.com/ionous/sashimi/script"
	"github.com/ionous/sashimi/standard"
	"github.com/stretchr/testify/assert"
	"testing"
)

//
func TestStandardRules(t *testing.T) {
	s := InitScripts()
	_, err := NewTestGame(t, s)
	assert.NoError(t, err)
}

//
func TestObjectSet(t *testing.T) {
	s := InitScripts()
	s.The("kinds",
		Have("amBlank", "text"),
		Have("amSet", "text"))

	s.The("kind",
		Called("test"),
		Has("amSet", "original"))

	g, err := NewTestGame(t, s)
	if assert.NoError(t, err) && assert.NotNil(t, g.Model) {
		if inst, ok := g.Model.Instances.FindInstance("test"); assert.True(t, ok) {
			gobj, exists := g.Game.Objects[inst.Id()]
			if assert.True(t, exists, "test instance should exist") && assert.NotNil(t, gobj) {
				obj := R.NewObjectAdapter(g.Game, gobj)
				assert.Equal(t, "original", obj.Text("amSet"), "should have original value")

				obj.SetText("amSet", "new")
				assert.Equal(t, "new", obj.Text("amSet"), "should change to new value")

				assert.Empty(t, obj.Text("amBlank"))

				obj.SetText("amBlank", "not blank any more")
				assert.NotEmpty(t, obj.Text("amBlank"))
			}
		}
	}
}

//
func TestStartupText(t *testing.T) {
	s := InitScripts()

	s.The("story",
		Called("testing"),
		Has("author", "me"),
		Has("headline", "extra extra"))

	s.The("room",
		Called("somewhere"),
		Has("description", "an empty room"),
		When("describing").Always(func(g G.Play) {
			g.StopHere()
		}),
	)

	if game, err := NewTestGame(t, s); assert.NoError(t, err, "compile should work") {
		story := game.FindFirstOf(game.Model.Classes.FindClass("stories"))
		assert.NotNil(t, story, "should have game")

		room := game.FindFirstOf(game.Model.Classes.FindClass("rooms"))
		assert.NotNil(t, room, "should have room")

		err = game.SendEvent("starting to play", story.String())
		assert.NoError(t, err, "starting to play")

		expected := []string{
			"", // FIX: this line shouldnt exist
			"testing",
			"extra extra by me",
			standard.VersionString,
			"",
			"somewhere",
			"an empty room",
		}
		assert.Exactly(t, expected, game.FlushOutput())
	}
}
