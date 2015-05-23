package sashimi

import (
	C "github.com/ionous/sashimi/console"
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
	"os"
	"testing"
)

//
func TestUnknownAction(t *testing.T) {
	s := Script{}
	s.The("kinds",
		When("this does not exists").Always(func(g G.Play) {
		}),
	)
	model, err := s.Compile(os.Stdout)
	if err == nil {
		model.PrintModel(t.Log)
		t.Fatal("expected failure")
	}
	t.Log("okay:", err)
}

//
func TestKnownAction(t *testing.T) {
	s := Script{}
	s.The("kinds",
		When("this exists").Always(func(g G.Play) {}),
		Can("exist").And("this exists").RequiresNothing())
	model, err := s.Compile(os.Stdout)
	if err != nil {
		t.Fatal(err)
	}
	model.PrintModel(t.Log)
}

//
func TestClassCallback(t *testing.T) {
	s := Script{}
	s.The("kinds",
		Have("description", "text"),
		Can("test").And("testing").RequiresNothing(),
		When("testing").Always(func(g G.Play) {
			t.Log("got the testing event")
			src := g.The("action.Source")
			if src.Name() != "obj" {
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
	cons := C.NewBufCon(nil)
	g, err := CompileGameWithConsole(&s, cons)
	if err != nil {
		t.Error("error:", err)
	}
	g.PushParserSource(func(g G.Play) G.IObject {
		return g.The("obj")
	})
	if err := g.SendEvent("testing", "obj"); err != nil {
		t.Error("error:", err)
	}
	if err := g.ProcessEvents(); err != nil {
		t.Error("error:", err)
	}
	out := cons.Flush()
	if len(out) != 1 || out[0] != "it's a trap!" {
		t.Fatal("mismatched output", out)
	}
}

//
func TestCallbackBeforeAfter(t *testing.T) {
	s := Script{}
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

	con := C.NewBufCon(nil)
	g, err := CompileGameWithConsole(&s, con)
	if err != nil {
		t.Error("error:", err)
	}
	g.PushParserSource(func(g G.Play) G.IObject {
		return g.The("obj")
	})
	if err := g.SendEvent("testing", "obj"); err != nil {
		t.Error("error:", err)
	}
	if err := g.ProcessEvents(); err != nil {
		t.Error("error:", err)
	}
	out := con.Flush()
	if len(out) != 2 || out[0] != "Before" || out[1] != "After" {
		t.Fatal("mismatched output", out)
	}
}

//
func TestCallbackParsing(t *testing.T) {
	s := Script{}

	s.The("kinds",
		Have("description", "text"),
		Can("test").And("testing").RequiresOne("kind"),
	)
	s.The("kinds",
		When("testing").Always(func(g G.Play) {
			t.Log("got testing")
			desc := g.The("action.Target").Text("description")
			t.Log("got desc", desc)
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
	strs := []string{"look at lookee"}
	con := C.NewBufCon(strs)
	g, err := CompileGameWithConsole(&s, con)
	if err != nil {
		t.Error(err)
	}
	numNames := len(g.Model.NounNames)
	if numNames != 2 {
		t.Error(numNames, "should equal", 2)
	}
	t.Log(g.Model.NounNames)

	g.PushParserSource(func(g G.Play) G.IObject {
		return g.The("looker")
	})
	g.RunForever()
	out := con.Flush()
	expect := []string{"look it's a test!"}
	if len(expect) != len(out) || expect[0] != out[0] {
		t.Fatal("Expected:", expect, "Actual:", out)
	}
}
