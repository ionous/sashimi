package tests

import (
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

//
func TestUnknownAction(t *testing.T) {
	s := &Script{}
	s.The("kinds",
		When("this does not exists").Always(func(g G.Play) {
		}),
	)
	if model, err := s.Compile(os.Stderr); !assert.Error(t, err, "expected failure") {
		model.PrintModel(t.Log)
	}
}

//
func TestKnownAction(t *testing.T) {
	s := &Script{}
	s.The("kinds",
		When("this exists").Always(func(g G.Play) {}),
		Can("exist").And("this exists").RequiresNothing())
	if model, err := s.Compile(os.Stderr); assert.NoError(t, err) {
		model.PrintModel(t.Log)
	}
}

//
func TestClassCallback(t *testing.T) {
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
	if g, err := NewTestGame(t, s); assert.NoError(t, err) {
		g.PushParserSource(func(g G.Play) G.IObject {
			return g.The("obj")
		})
		if err := g.SendEvent("testing", "Obj"); assert.NoError(t, err) {
			if err := g.ProcessEvents(); assert.NoError(t, err) {
				expected := []string{"it's a trap!"}
				assert.EqualValues(t, expected, g.FlushOutput())
			}
		}
	}
}

//
func TestCallbackBeforeAfter(t *testing.T) {
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
	if g, err := NewTestGame(t, s); assert.NoError(t, err) {
		g.PushParserSource(func(g G.Play) G.IObject {
			return g.The("obj")
		})
		if err := g.SendEvent("testing", "Obj"); assert.NoError(t, err) {
			if err := g.ProcessEvents(); assert.NoError(t, err) {
				expected := []string{"Before", "After"}
				assert.EqualValues(t, expected, g.FlushOutput())
			}
		}
	}
}

//
func TestCallbackParsing(t *testing.T) {
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
	if g, err := NewTestGame(t, s); assert.NoError(t, err) {
		if assert.Len(t, g.Model.NounNames, 2) {
			g.PushParserSource(func(g G.Play) G.IObject {
				return g.The("looker")
			})
			str := "look at lookee"
			expected := []string{"look it's a test!"}
			assert.EqualValues(t, expected, g.RunInput(str).FlushOutput())
		}
	}
}
