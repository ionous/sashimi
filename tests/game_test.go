package tests

import (
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
	"github.com/ionous/sashimi/standard/live"
	"github.com/ionous/sashimi/util/ident"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	test, err := NewTestGame(t, s)
	if assert.NoError(t, err) && assert.NotNil(t, test.Model) {
		if inst, ok := test.Model.Instances[ident.MakeId("test")]; assert.True(t, ok) {
			gobj, exists := test.Game.GetInstance(inst.Id)
			if assert.True(t, exists, "test instance should exist") && assert.NotNil(t, gobj) {
				adapter := test.Game.NewAdapter()
				obj := adapter.NewGameObject(gobj)
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
	//R.DebugGet = true

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

	bannerCalled := false
	storyExists := false
	nameOkay := false
	s.The("stories",
		When("printing the banner").Always(func(g G.Play) {
			bannerCalled = true
			story := g.The("story")
			storyExists = story.Exists()
			name := story.Text("name")
			nameOkay = name == "testing"
		}))

	if test, err := NewTestGame(t, s); assert.NoError(t, err, "compile should work") {
		if out, err := test.Commence(); assert.NoError(t, err, "couldnt commence") {

			expected := []string{
				"testing.",
				"extra extra by me.",
				live.VersionString,
				"",
				"somewhere",
				"an empty room",
				"",
			}
			if assert.True(t, bannerCalled, "banner called") {
				require.True(t, storyExists, "story exists")
				require.True(t, nameOkay, "name set")
			}
			require.Exactly(t, expected, out)
		}
	}
}
