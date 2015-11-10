package tests

import (
	G "github.com/ionous/sashimi/game"
	R "github.com/ionous/sashimi/runtime"
	"github.com/ionous/sashimi/runtime/api"
	. "github.com/ionous/sashimi/script"
	"github.com/ionous/sashimi/standard"
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
		if inst, ok := test.Model.Instances.FindInstance("test"); assert.True(t, ok) {
			gobj, exists := test.Game.ModelApi.GetInstance(inst.Id)
			if assert.True(t, exists, "test instance should exist") && assert.NotNil(t, gobj) {
				obj := R.NewObjectAdapter(test.Game, gobj)
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
	R.DebugGet = true

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
		if story, ok := api.FindFirstOf(test.Game.ModelApi, ident.MakeId("stories")); assert.True(t, ok, "should have test story") {
			if _, ok := api.FindFirstOf(test.Game.ModelApi, ident.MakeId("rooms")); assert.True(t, ok, "should have room") {
				err = test.Game.QueueEvent("commencing", story.GetId())
				require.NoError(t, err, "commencing")

				expected := []string{
					"testing",
					"extra extra by me",
					standard.VersionString,
					"",
					"somewhere",
					"an empty room",
					"",
				}
				if out, e := test.FlushOutput(); assert.NoError(t, e) {
					if assert.True(t, bannerCalled, "banner called") {
						require.True(t, storyExists, "story exists")
						require.True(t, nameOkay, "name set")
					}
					require.Exactly(t, expected, out)
				}
			}
		}
	}
}
