package tests

import (
	"github.com/ionous/mars/core"
	"github.com/ionous/mars/g"
	"github.com/ionous/mars/rt"
	. "github.com/ionous/mars/script"
	"github.com/stretchr/testify/assert"
	"testing"
)

func lines(v ...string) []string {
	// TBD, not sure where this trailing blank line is  coming from...
	return append(v, "")
}

func TestRawTextProperty(t *testing.T) {
	script := &Script{
		The("kinds", Called("actors"), Have("greeting", "text")),
		The("actor", Called("player"), Has("greeting", "hello world")),
	}

	if test, err := NewTestGame(t, script); assert.NoError(t, err, "new game") {
		if player, ok := test.Metal.GetInstance("player"); assert.True(t, ok, "found world") {
			if greeting, ok := player.FindProperty("greeting"); assert.True(t, ok, "has greeting") {
				g := greeting.GetGeneric()
				if v, ok := g.(rt.TextEval); assert.True(t, ok, "text eval") {
					run := test.Game.Rtm
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
	script := &Script{
		The("kinds", Called("actors"), Have("counter", "num")),
		The("actor", Called("player"), Has("counter", core.AddNum{core.N(2), core.N(3)})),
	}
	if test, err := NewTestGame(t, script); assert.NoError(t, err, "new game") {
		if player, ok := test.Metal.GetInstance("player"); assert.True(t, ok, "found player") {
			if counter, ok := player.FindProperty("counter"); assert.True(t, ok, "has greeting") {
				g := counter.GetGeneric()
				if v, ok := g.(rt.NumEval); assert.True(t, ok, "num eval") {
					run := test.Game.Rtm
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

func TestActionNames(t *testing.T) {
	script := &Script{
		The("kinds", Called("actors"), Have("greeting", "text")),
		The("actor", Called("player"), Has("greeting", "hello world")),
		The("kinds", Called("actors"),
			Can("greet the world").And("greeting the world").RequiresNothing(),
			To("greet the world",
				g.Say(g.The("player").Text("greeting")),
				g.Say(g.The("action.source").Text("greeting")),
				g.Say(g.The("actor").Text("greeting")),
			)),
	}
	//running queued action
	//got changed value hello
	if test, err := NewTestGame(t, script); assert.NoError(t, err, "new game") {
		if err := test.Game.RunAction(core.MakeStringId("greet the world"), g.The("player")); assert.NoError(t, err, "run action") {
			if v, err := test.FlushOutput(); assert.NoError(t, err, "process") {
				expected := lines("hello world",
					"hello world",
					"hello world")
				if !assert.EqualValues(t, expected, v) {
					t.FailNow()
				}
			}
		}
	}
}

func TestTarget(t *testing.T) {
	script := &Script{
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
		if err := test.Game.RunAction(core.MakeStringId("greet actor"), g.The("player"), g.The("npc")); assert.NoError(t, err, "run action") {
			if v, err := test.FlushOutput(); assert.NoError(t, err, "process") {
				if !assert.EqualValues(t, lines("hello npc"), v) {
					t.FailNow()
				}
			}
		}
	}
}

// TestRun calls an action from an action
func TestRun(t *testing.T) {
	script := &Script{
		The("kinds", Called("actors"), Have("greeting", "text")),
		The("actor", Called("player"), Has("greeting", "hello")),
		The("actor", Called("npc"), Exists()),
		The("kinds", Called("actors"),
			Can("test nothing").And("testing nothing").RequiresNothing(),
			To("test nothing",
				g.Say("absolutely nothing")),
			Can("greet actor").And("greeting actor").RequiresOne("actor"),
			To("greet actor",
				g.The("player").Go("test nothing"),
			)),
	}
	if test, err := NewTestGame(t, script); assert.NoError(t, err, "new game") {
		if err := test.Game.RunAction(core.MakeStringId("greet actor"), g.The("player"), g.The("npc")); assert.NoError(t, err, "run action") {
			if v, err := test.FlushOutput(); assert.NoError(t, err, "process") {
				if !assert.EqualValues(t, lines("absolutely nothing"), v) {
					t.FailNow()
				}
			}
		}
	}
}

func TestStopHere(t *testing.T) {
	script := &Script{
		The("kinds", Called("actors"), Have("greeting", "text")),
		The("actor", Called("player"), Has("greeting", "hello world")),
		The("kinds", Called("actors"),
			Can("greet the world").And("greeting the world").RequiresNothing(),
			When("greeting the world").Always(
				g.Say(g.The("player").Text("greeting")),
				g.StopHere(),
				g.Say(g.The("action.source").Text("greeting")),
				g.Say(g.The("actor").Text("greeting")),
			),
		),
	}
	//running queued action
	//got changed value hello
	if test, err := NewTestGame(t, script); assert.NoError(t, err, "new game") {
		if err := test.Game.RunAction(core.MakeStringId("greet the world"), g.The("player")); assert.NoError(t, err, "run action") {
			if v, err := test.FlushOutput(); assert.NoError(t, err, "process") {
				if !assert.EqualValues(t,
					lines("hello world"), v) {
					t.FailNow()
				}
			}
		}
	}
}
