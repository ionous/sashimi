package tests

import (
	M "github.com/ionous/sashimi/model"
	R "github.com/ionous/sashimi/runtime"
	. "github.com/ionous/sashimi/script"
	"github.com/ionous/sashimi/standard"
	"github.com/ionous/sashimi/util/ident"
	"github.com/stretchr/testify/assert"

	"os"
	"testing"
)

func makeSweets(s *Script) {
	s.The("kinds",
		Called("sweets"),
		Have("desc", "text"),
		Have("price", "num"),
		AreOneOf("delicious", "decent", "acceptable", "you can't be serious"))
}
func nameSweets(s *Script) {
	s.The("sweets",
		Table("name", "desc", "delicious-property").Contains(
			"Boreo", "Creme filled wafer things.", "acceptable").And(
			"Vegan chocolate chip cookies", "The secret is the bacon... err... baking.", "delicious").And(
			"Sugar coated ants", "A crunchy summer's day snack.", "you can't be serious",
		),
	)
}

// 1. create a simple table declaration and generate some fixed instances
func TestTableDecl(t *testing.T) {
	s := &Script{}
	makeSweets(s)
	nameSweets(s)
	if m, err := s.Compile(os.Stderr); assert.NoError(t, err, "table compile") {
		assert.Len(t, m.Instances, 3)
		// and now some values:
		if inst, ok := m.Instances.FindInstance("boreo"); assert.True(t, ok, "find tabled instance by name") {
			//
			if val, ok := inst.ValueByName("desc"); assert.True(t, ok, "find desc") {
				if str, ok := val.(string); assert.True(t, ok, "have string") {
					assert.EqualValues(t, str, "Creme filled wafer things.")
				}
			}
			//
			if val, ok := inst.ValueByName("delicious-property"); assert.True(t, ok, "find deliciousness") {
				assert.EqualValues(t, val, M.MakeStringId("acceptable"))
			}
		}
	}
}

// 2. create a simple table declaration with autogenerated instances
func TestTableGeneration(t *testing.T) {
	s := &Script{}
	makeSweets(s)
	s.The("sweets",
		Table("desc", "delicious-property").Contains(
			"Creme filled wafer things.", "acceptable").And(
			"It looks the way poisonous berries smell.", "delicious").And(
			"A crunchy summer's day snack.", "you can't be serious",
		),
	)
	if m, err := s.Compile(os.Stderr); assert.NoError(t, err, "table compile") {
		assert.Len(t, m.Instances, 3)
	}
}

// 3. merge some non-contrary instance data
func TestTabledData(t *testing.T) {
	s := &Script{}
	makeSweets(s)
	nameSweets(s)
	s.Our("Boreo", Is("acceptable"), Has("price", 42))
	if m, err := s.Compile(os.Stderr); assert.NoError(t, err, "table compile") {
		if inst, ok := m.Instances.FindInstance("boreo"); assert.True(t, ok, "find tabled instance by name") {
			if val, ok := inst.ValueByName("price"); assert.True(t, ok, "find desc") {
				assert.EqualValues(t, 42, val)
			}
		}
	}
}

func makePeople(s *Script) *Script {
	s.The("kinds", Called("objects"), Exist())
	s.The("objects",
		Called("actors"),
		Have("favorite sweet", "sweet"))
	return s
}

func namePeople(s *Script) {
	s.The("actors",
		Table("name", "favorite sweet").Contains(
			"Marvin", "Sugar coated ants").And(
			"Allen", "Boreo").And(
			"Grace", "Vegan chocolate chip cookies",
		),
	)
}

// 4. test table declarations which use pointers
func TestTablePointers(t *testing.T) {
	s := &Script{}
	makeSweets(s)
	nameSweets(s)
	makePeople(s)
	namePeople(s)
	if m, err := s.Compile(os.Stderr); assert.NoError(t, err, "table compile") {
		if inst, ok := m.Instances.FindInstance("Grace"); assert.True(t, ok, "find person by name") {
			if val, ok := inst.ValueByName("Favorite Sweet"); assert.True(t, ok, "find favorite") {
				if id, ok := val.(ident.Id); assert.True(t, ok, "id") {
					assert.EqualValues(t, id, "VeganChocolateChipCookies")
				}
			}
		}
	}
}

// 5. test runtime instance usage which use pointers
func TestTableRuntime(t *testing.T) {
	s := standard.InitStandardLibrary()
	makeSweets(s)
	nameSweets(s)
	makePeople(s)
	namePeople(s)
	if g, err := NewTestGame(t, s); assert.NoError(t, err) {
		if gobj, ok := g.FindObject("Grace"); assert.True(t, ok) {
			grace := R.NewObjectAdapter(g.Game, gobj)
			if sweet := grace.Object("favorite sweet"); assert.True(t, sweet.Exists(), "grace should have a treat") {
				if gobj, ok := g.FindObject("Boreo"); assert.True(t, ok) {
					boreo := R.NewObjectAdapter(g.Game, gobj)
					grace.SetObject("favorite sweet", boreo)
					retread := grace.Object("favorite sweet").Name()
					assert.Equal(t, retread, "Boreo")
				}
			}
		}
	}
}
