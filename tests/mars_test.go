package tests

import (
	"github.com/ionous/mars/core"
	"github.com/ionous/mars/g"
	"github.com/ionous/mars/rt"
	"github.com/ionous/mars/rtm"
	"github.com/ionous/mars/script"
	. "github.com/ionous/mars/script/s"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRawTextProperty(t *testing.T) {
	script := &script.Script{
		The("kinds", Called("actors"), Have("greeting", "text")),
		The("actor", Called("player"), Has("greeting", "hello world")),
	}

	// FIX FIX FIX FIX FIX -- the test game -- any game -- shouldnt require a parser.
	// that should be on the front end, wrapping the game.
	// ditto the "player"
	// the understandings used by the parser can just sit there
	// in the future, maybe we could put the understanding in an outer layer
	if test, err := NewTestGame(t, script); assert.NoError(t, err, "new game") {
		if player, ok := test.Game.GetInstance("player"); assert.True(t, ok, "found world") {
			if greeting, ok := player.FindProperty("greeting"); assert.True(t, ok, "has greeting") {
				g := greeting.GetGeneric()
				if v, ok := g.(rt.TextEval); assert.True(t, ok, "text eval") {
					run := rtm.NewRtm(test.Game.Model)
					if text, e := v.GetText(run); assert.NoError(t, e, "got text") {
						if !assert.Equal(t, "hello world", text.String()) {
							t.FailNow()
						}
					}
				}
			}
		}
	}
}

func TestNumEvalProperty(t *testing.T) {
	script := &script.Script{
		The("kinds", Called("actors"), Have("counter", "num")),
		The("actor", Called("player"), Has("counter", core.AddNum{core.N(2), core.N(3)})),
	}

	// FIX FIX FIX FIX FIX -- the test game -- any game -- shouldnt require a parser.
	// that should be on the front end, wrapping the game.
	// ditto the "player"
	// the understandings used by the parser can just sit there
	// in the future, maybe we could put the understanding in an outer layer
	if test, err := NewTestGame(t, script); assert.NoError(t, err, "new game") {
		if player, ok := test.Game.GetInstance("player"); assert.True(t, ok, "found player") {
			if counter, ok := player.FindProperty("counter"); assert.True(t, ok, "has greeting") {
				g := counter.GetGeneric()
				if v, ok := g.(rt.NumEval); assert.True(t, ok, "num eval") {
					run := rtm.NewRtm(test.Game.Model)
					if num, e := v.GetNumber(run); assert.NoError(t, e, "got num") {
						if !assert.EqualValues(t, 5, num.Float()) {
							t.FailNow()
						}
					}
				}
			}
		}
	}
}

func TestOldStyleAction(t *testing.T) {
	script := &script.Script{
		The("kinds", Called("actors"), Have("greeting", "text")),
		The("actor", Called("player"), Has("greeting", "hello world")),
		The("kinds", Called("actors"),
			Can("greet the world").And("greeting the world").RequiresNothing(),
			To("greet the world",
				g.Say(g.The("player").Text("greeting")),
			)),
	}
	//running queued action
	//got changed value hello
	if test, err := NewTestGame(t, script); assert.NoError(t, err, "new game") {
		if _, err := test.Game.QueueAction("greet the world", "player"); assert.NoError(t, err, "queue") {
			if v, err := test.FlushOutput(); assert.NoError(t, err, "process") {
				if !assert.EqualValues(t, "hello world", v[0]) {
					t.FailNow()
				}
			}
		}
	}
}

func TestOldStyleTarget(t *testing.T) {
	script := &script.Script{
		The("kinds", Called("actors"), Have("greeting", "text")),
		The("actor", Called("player"), Has("greeting", "hello")),
		The("actor", Called("npc"), Exists()),
		The("kinds", Called("actors"),
			Can("greet actor").And("greeting actor").RequiresOne("actor"),
			To("greet actor",
				g.Say(g.The("player").Text("greeting"), g.The("action.target").Text("name")),
			)),
	}
	if test, err := NewTestGame(t, script); assert.NoError(t, err, "new game") {
		if _, err := test.Game.QueueAction("greet actor", "player", "npc"); assert.NoError(t, err, "queue") {
			if v, err := test.FlushOutput(); assert.NoError(t, err, "process") {
				if !assert.EqualValues(t, "hello npc", v[0]) {
					t.FailNow()
				}
			}
		}
	}
}
