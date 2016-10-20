package tests

import (
	"github.com/ionous/mars/core"
	"github.com/ionous/mars/rt"
	"github.com/ionous/mars/rtm"
	"github.com/ionous/mars/script"
	. "github.com/ionous/mars/script/s"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRawProperty(t *testing.T) {
	script := &script.Script{
		The("kinds",
			Called("rooms"),
			Have("greeting", "text"),
		),
		The("kinds", Called("actors"), Exist()),
		The("actor", Called("player"), Exists()),
		The("room", Called("world"), Has("greeting", "hello")),
	}

	// FIX FIX FIX FIX FIX -- the test game -- any game -- shouldnt require a parser.
	// that should be on the front end, wrapping the game.
	// ditto the "player"
	// the understandings used by the parser can just sit there
	// in the future, maybe we could put the understanding in an outer layer
	if test, err := NewTestGame(t, script); assert.NoError(t, err, "created game") {
		//
		if world, ok := test.Game.GetInstance("world"); assert.True(t, ok, "found world") {
			//
			if greeting, ok := world.FindProperty("greeting"); assert.True(t, ok, "has greeting") {
				//
				g := greeting.GetGeneric()
				if v, ok := g.(rt.TextEval); assert.True(t, ok, "text eval") {
					//
					run := rtm.NewRtm(test.Game.Model, nil)
					if text, e := v.GetText(run); assert.NoError(t, e, "got text") {
						assert.Equal(t, "hello", text.String())
						t.Log(text)
					}
				}
			}
		}
	}
}
func TestNumEvalProperty(t *testing.T) {
	script := &script.Script{
		The("kinds",
			Called("actors"),
			Have("counter", "num"),
		),
		The("actor", Called("player"), Has("counter", core.AddNum{core.N(2), core.N(3)})),
	}

	// FIX FIX FIX FIX FIX -- the test game -- any game -- shouldnt require a parser.
	// that should be on the front end, wrapping the game.
	// ditto the "player"
	// the understandings used by the parser can just sit there
	// in the future, maybe we could put the understanding in an outer layer
	if test, err := NewTestGame(t, script); assert.NoError(t, err, "created game") {
		if player, ok := test.Game.GetInstance("player"); assert.True(t, ok, "found player") {
			if counter, ok := player.FindProperty("counter"); assert.True(t, ok, "has greeting") {
				g := counter.GetGeneric()
				if v, ok := g.(rt.NumEval); assert.True(t, ok, "num eval") {
					run := rtm.NewRtm(test.Game.Model, nil)
					if num, e := v.GetNumber(run); assert.NoError(t, e, "got num") {
						assert.EqualValues(t, 5, num.Float())
					}
				}
			}
		}
	}
}
