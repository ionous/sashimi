package tests

import (
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
	"github.com/ionous/sashimi/util/ident"
	"github.com/stretchr/testify/assert"
	"testing"
)

//
func TestActionUnknown(t *testing.T) {
	s := &Script{}
	s.The("kinds",
		When("this does not exists").Always(func(g G.Play) {
		}),
	)
	if model, err := s.Compile(Log(t)); !assert.Error(t, err, "expected failure") {
		model.Model.PrintModel(t.Log)
	}
}

//
func TestActionKnown(t *testing.T) {
	s := &Script{}
	s.The("kinds",
		When("this exists").Always(func(g G.Play) {}),
		Can("exist").And("this exists").RequiresNothing())
	if model, err := s.Compile(Log(t)); assert.NoError(t, err) {
		model.Model.PrintModel(t.Log)
	}
}

//
func TestActionClassCallback(t *testing.T) {
	s := &Script{}
	s.The("kinds",
		Have("description", "text"),
		Can("test").And("testing").RequiresNothing(),
		When("testing").Always(func(g G.Play) {
			t.Log("got the testing event")
			src := g.The("action.Source")
			if src.Text("Name") != "obj" {
				t.Error("wrong object", src)
			}
			desc := src.Text("description")
			if desc != "it's a trap!" {
				t.Error("wrong desc", desc)
			}
			g.Say(desc)
		}),
	)
	s.The("kind",
		Called("obj"),
		Has("description", "it's a trap!"),
	)
	s.The("kind",
		Called("other"),
		Has("description", "it's an error!"),
	)
	if g, err := NewTestGameSource(t, s, "obj", nil); assert.NoError(t, err) {
		if _, err := g.Game.QueueAction("test", ident.MakeId("obj")); assert.NoError(t, err) {
			if err := g.Game.ProcessActions(); assert.NoError(t, err) {
				if out, err := g.FlushOutput(); assert.NoError(t, err) {
					expected := []string{"it's a trap!"}
					assert.EqualValues(t, expected, out)
				}
			}
		}
	}
}

//
func TestActionCallbackBeforeAfter(t *testing.T) {
	s := &Script{}
	s.The("kinds",
		Can("test").And("testing").RequiresNothing(),
		When("testing").Always(func(g G.Play) {
			g.Say("After")
		}),
		Before("testing").Always(func(g G.Play) {
			g.Say("Before")
		}),
	)
	s.The("kind", Called("obj"), Exists())
	if g, err := NewTestGameSource(t, s, "obj", nil); assert.NoError(t, err) {
		if _, err := g.Game.QueueAction("test", ident.MakeId("obj")); assert.NoError(t, err) {
			if err := g.Game.ProcessActions(); assert.NoError(t, err) {
				if out, err := g.FlushOutput(); assert.NoError(t, err) {
					expected := []string{"Before", "After"}
					assert.EqualValues(t, expected, out)
				}
			}
		}
	}
}

//
func TestActionCallbackParsing(t *testing.T) {
	s := &Script{}

	s.The("kinds",
		Have("description", "text"),
		Can("test").And("testing").RequiresOne("kind"),
	)
	s.The("kinds",
		When("testing").Always(func(g G.Play) {
			desc := g.The("action.Target").Text("description")
			g.Say(desc)
		}),
	)
	s.The("kind",
		Called("looker"),
		Exists(),
	)
	s.The("kind",
		Called("lookee"),
		Has("description", "look it's a test!"),
	)
	s.Execute("test",
		Matching("look|l at {{something}}"),
	)
	// should trigger "test", which should print the description
	if g, err := NewTestGameSource(t, s, "looker", nil); assert.NoError(t, err) {
		if assert.Len(t, g.Model.Aliases, 2) {
			str := "look at lookee"
			if res, err := g.RunInput(str); assert.NoError(t, err, "handle input") {
				expected := []string{"look it's a test!"}
				assert.EqualValues(t, expected, res)
			}
		}
	}
}
