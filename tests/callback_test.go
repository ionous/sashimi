package tests

import (
	. "github.com/ionous/mars/core"
	. "github.com/ionous/mars/script"
	"github.com/ionous/mars/script/g"
	"github.com/ionous/sashimi/compiler"
	S "github.com/ionous/sashimi/source"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

// TestCallbackUnknown tests compiler failure when an action does not exist
func TestCallbackUnknown(t *testing.T) {
	s := NewScript(
		The("kinds", When("this does not exists").Always(DoNothing{})),
	)
	src := &S.Statements{}
	if err := s.Generate(src); assert.NoError(t, err, "build") {
		if _, err := compiler.Compile(*src, ioutil.Discard); assert.Error(t, err, "expected failure") {
			return
		}
	}
	t.FailNow()
}

//TestCallbackKnown tests compiler success for a simple action
func TestCallbackKnown(t *testing.T) {
	s := NewScript(
		The("kinds",
			When("this exists").Always(DoNothing{}),
			Can("exist").And("this exists").RequiresNothing()),
	)
	//"couldnt compile action ### couldn't find class "
	src := &S.Statements{}
	if err := s.Generate(src); assert.NoError(t, err, "build") {
		if _, err := compiler.Compile(*src, ioutil.Discard); assert.NoError(t, err, "expected success") {
			return
		}
	}
	t.FailNow()
}

// TestCallbackClass tests the execution of a simple callback
func TestCallbackClass(t *testing.T) {
	s := NewScript(
		The("kinds",
			Have("description", "text"),
			Can("test").And("testing").RequiresNothing(),
			When("testing").Always(
				Choose{
					If: IsText{
						g.The("action.Source").Text("Name"),
						NotEqual,
						T("obj")},
					True: Error{T("wrong object")},
				},
				Choose{
					If: IsText{
						g.The("action.Source").Text("description"),
						NotEqual,
						T("it's a trap!")},
					True:  Error{T("wrong description")},
					False: g.Say(g.The("action.Source").Text("description")),
				},
			),
		),
		The("kind",
			Called("obj"),
			Has("description", "it's a trap!")),
		The("kind",
			Called("other"),
			Has("description", "it's an error!")),
	)
	if test, err := NewTestGameScript(t, s, "obj", nil); assert.NoError(t, err) {
		if err := test.RunNamedAction("test", g.The("obj")); assert.NoError(t, err) {
			if out, err := test.FlushOutput(); assert.NoError(t, err) {
				expected := lines("it's a trap!")
				if assert.EqualValues(t, expected, out) {
					return
				}
			}
		}
	}
	t.FailNow()
}

// TestCallbackBeforeAfter: capture actions before and after an event.
func TestCallbackBeforeAfter(t *testing.T) {
	s := NewScript(
		The("kinds",
			Can("test").And("testing").RequiresNothing(),
			When("testing").Always(g.Say("After")),
			Before("testing").Always(g.Say("Before")),
		),
		The("kind", Called("obj"), Exists()),
	)
	if test, err := NewTestGameScript(t, s, "obj", nil); assert.NoError(t, err) {
		if err := test.RunNamedAction("test", g.The("obj")); assert.NoError(t, err) {
			if out, err := test.FlushOutput(); assert.NoError(t, err) {
				expected := lines("Before", "After")
				if assert.EqualValues(t, expected, out) {
					return
				}
			}
		}
	}
	t.FailNow()
}
