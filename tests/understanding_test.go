package tests

import (
	"github.com/ionous/mars/g"
	. "github.com/ionous/mars/script"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestUnderstanding should trigger "test", which should print a description.
func TestUUManually(t *testing.T) {
	s := understandingTest()
	expected := lines("look it's a test!")
	//
	if test, err := NewTestGameSource(t, s, "looker", nil); assert.NoError(t, err) {
		if err := test.Game.RunAction("test", g.The("looker"), g.The("lookee")); assert.NoError(t, err, "run manually") {
			if res, err := test.FlushOutput(); assert.NoError(t, err, "raw flush") {
				if assert.EqualValues(t, expected, res, "raw output") {
					return
				}
			}
		}
	}
	t.FailNow()
}

func TestUParserText(t *testing.T) {
	s := understandingTest()
	expected := lines("look it's a test!")
	//
	if test, err := NewTestGameSource(t, s, "looker", nil); assert.NoError(t, err) {
		if assert.Len(t, test.Metal.Model.Aliases, 2, "parsed actions") {
			str := "look at lookee"
			if res, err := test.RunInput(str); assert.NoError(t, err, "handle input") {
				if assert.EqualValues(t, expected, res, "parsed output") {
					return
				}
			}
		}
	}
	t.FailNow()
}

func TestUKnownAs(t *testing.T) {
	expected := lines("look it's a test!")
	s := append(*understandingTest(),
		The("lookee", IsKnownAs("something special")))
	//
	if test, err := NewTestGameSource(t, &s, "looker", nil); assert.NoError(t, err) {
		ok := "look at something special"
		if res, err := test.RunInput(ok); assert.NoError(t, err, "something special") {
			if assert.EqualValues(t, expected, res, "parsed output") {
				ng := "look at nothing special"
				if _, err := test.RunInput(ng); assert.Error(t, err, "nothing special") {
					return
				}
			}
		}
	}
	t.FailNow()
}

func understandingTest() *Script {
	return &Script{
		The("kinds",
			Have("description", "text"),
			Can("test").And("testing").RequiresOne("kind"),
		),
		The("kinds",
			When("testing").Always(
				g.Say(g.The("action.Target").Text("description")),
			),
		),
		The("kind",
			Called("looker"),
			Exists(),
		),
		The("kind",
			Called("lookee"),
			Has("description", "look it's a test!"),
		),
		Understand("look|l at {{something}}").As("test"),
	}
}
