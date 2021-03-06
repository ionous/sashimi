package tests

import (
	"fmt"
	"github.com/ionous/sashimi/meta"
	. "github.com/ionous/sashimi/script"
	"github.com/ionous/sashimi/standard"
	"github.com/ionous/sashimi/util/ident"
	"github.com/stretchr/testify/assert"
	"testing"
)

func makeTestRoom() *Script {
	s := standard.InitStandardLibrary()
	// //exit door and its room, with optional door
	s.The("room", Called("the lobby"), Exists(),
		// two-way direction
		Going("up").Through("the trap door").ConnectsTo("the parapet"),
		// one-way directions
		Going("down").ArrivesAt("the basement"),
	)
	s.The("foyer",
		// direction to room, reverses
		Going("north").ConnectsTo("the outside"),
		// direction to room, no-reverse
		Going("west").ArrivesAt("the lobby"),
	)
	s.The("lobby",
		// non-commensurate directions
		Going("north").ArrivesAt("the foyer"))
	// explicitly declaring the door should be fine.
	s.The("door", Called("the cellar door"), Exists())
	// direction to door, does not reverse
	s.The("basement", Going("south").
		ArrivesAt("the outside").Door("the cellar door"),
	)
	// not explicitly declaring the door should also work:
	//     The("door", Called("the cellar door"), Exists())
	// door-to-door two-way.
	s.The("foyer", Through("the curtain").
		ConnectsTo("the cloakroom").Door("the cloakroom-curtain"),
	)
	// FIX: want to map "name" to a property, and if it doesn't exist fall back on split id.
	// FIX? wonder if you could inject a report of some kind to pull in the description /brief of a door
	// automatically to match it's other side.
	s.The("door", Called("curtain"), Has("brief", "A red velvet curtain."))
	s.The("door", Called("cloakroom-curtain"), Has("brief", "A red velvet curtain."))
	return s
}

//
// test the creation of a connected world
//
func TestMoveConstruction(t *testing.T) {
	s := makeTestRoom()
	m, err := s.Compile(Log(t))
	if assert.NoError(t, err) {
		m.Model.PrintModel(t.Log)
	}
}

//
// test moving around
//
func TestMoveGoing(t *testing.T) {
	s := makeTestRoom()
	s.The("player", Exists(), In("the foyer"))
	if test, err := NewTestGame(t, s); assert.NoError(t, err) {
		if e := testMoves(t, test,
			xMove{"go west", "Lobby"},
			xMove{"go east", "Lobby"},
			xMove{"go up", "Parapet"},
			xMove{"go down", "Lobby"}, // first two way
			xMove{"go down", "Basement"},
			xMove{"go up", "Basement"},
			xMove{"go south", "Outside"},
			xMove{"go south", "Foyer"},
			xMove{"enter curtain", "Cloakroom"},
		); !assert.NoError(t, e, "failed move") {
			p := test.Parser
			t.Logf("parser has %d comprehension", len(p.Comprehensions))
			for k, v := range p.Comprehensions {
				t.Logf("%v:%v", k, v)
			}
		}
	}
}

type xMove struct {
	cmd string
	res string
}

func testMoves(t *testing.T, test TestGame, moves ...xMove) (err error) {
	// FIX: relations are stored in the model
	if p, ok := test.Game.GetInstance(ident.MakeId("player")); assert.True(t, ok, "found player") {
		for _, move := range moves {
			t.Logf("%s => %s", move.cmd, move.res)
			if out, e := test.RunInput(move.cmd); e != nil {
				err = e
				break
			} else if res := where(test.Game, p); res.String() != move.res {
				err = fmt.Errorf("unexpected move result: %v(expected) != %v(actual); %v", move.res, res, out)
				break
			}
		}
	}
	return err
}

func where(mdl meta.Model, gobj meta.Instance) (ret ident.Id) {
	if prop, ok := gobj.FindProperty("whereabouts"); ok {
		ret = prop.GetValue().GetObject()
	}
	return ret
}
