package tests

import (
	"encoding/json"
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/meta"
	R "github.com/ionous/sashimi/runtime"
	"github.com/ionous/sashimi/runtime/api"
	. "github.com/ionous/sashimi/script"
	"github.com/ionous/sashimi/standard/framework"
	. "github.com/ionous/sashimi/standard/live"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

//
func TestInvSaveLoad(t *testing.T) {
	s := InvScript()
	//
	s.The("player", Possesses("the candy"))
	s.Our("test story", When("commencing").Always(func(g G.Play) {
		g.Go(Give("the hat").To("the player"),
			Give("the hook").To("the player"),
			Move("the candy").OutOfWorld())
	}))
	//
	pc := func(mdl meta.Model) api.LookupParents {
		return framework.NewParentLookup(mdl)
	}
	if test, err := NewTestGameSource(t, s, "player", pc); assert.NoError(t, err) {
		//bytes, _ := json.MarshalIndent(test.Model.Instances, "", " ")
		//fmt.Println(string(bytes))
		g := test.Game.NewAdapter()
		candy := g.GetObject("candy")

		// candy's parent should be the player
		if parent, rel := candy.ParentRelation(); assert.True(t, parent.Exists(), "parent exists") {
			if assert.Equal(t, "owner", rel) {
				t.Log(candy, parent, rel)
			}
		}

		if _, err := test.Commence(); assert.NoError(t, err, "couldnt commence") {
			if _, err := R.SaveGame(g, false); assert.NoError(t, err, "couldnt save game") {
				m := make(map[string]string)
				if err := json.Unmarshal(test.saver.blob, &m); assert.NoError(t, err, "couldnt read blob") {
					require.Equal(t, m["Candy.ObjectsOwner"], "")
					require.Equal(t, m["Hat.ObjectsOwner"], "player")
					require.Equal(t, m["Hook.ObjectsOwner"], "player")
					require.Equal(t, m["Hook.ObjectsWhereabouts"], "")
					require.Equal(t, m["Neverland.RoomsVisited"], "visited")
					require.Equal(t, m["Player.ObjectsWhereabouts"], "neverland")
				}
			}
		}
	}
}

// TalkScript creates some dialog to test.
func InvScript() *Script {
	s := InitScripts()
	s.The("room", Called("neverland"), Exists())
	s.The("story", Called("test story"), Exists())
	s.The("prop", Called("hook"), Exists(), In("neverland"))
	s.The("prop", Called("candy"), Exists())
	s.The("prop", Called("hat"), Exists())
	return s
}
