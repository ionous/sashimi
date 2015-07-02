package sashimi

import (
	G "github.com/ionous/sashimi/game"
	M "github.com/ionous/sashimi/model"
	. "github.com/ionous/sashimi/script"
	"github.com/ionous/sashimi/standard"
	"github.com/stretchr/testify/assert"
	"os"
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
	s.The("door", Called("curtain"), Has("brief", "A red velvet curtain"))
	s.The("door", Called("cloakroom-curtain"), Has("brief", "A red velvet curtain"))
	return s
}

//
// test the creation of a connected world
//
func TestMoveConstruction(t *testing.T) {
	s := makeTestRoom()
	m, err := s.Compile(os.Stderr)
	if assert.NoError(t, err) {
		m.PrintModel(t.Log)
	}
}

//
// test moving around
//
func TestMoveGoing(t *testing.T) {
	s := makeTestRoom()
	s.The("player", Exists(), In("the foyer"))
	if g, err := NewTestGame(t, s); assert.NoError(t, err) {
		// FIX: move parser source into the model/parser
		g.PushParserSource(func(g G.Play) (ret G.IObject) {
			return g.The("player")
		})
		g.PushParentLookup(func(g G.Play, o G.IObject) (ret G.IObject) {
			if parent, where := standard.DirectParent(o); where != "" {
				ret = parent
			}
			return ret
		})
		testMoves(t, g,
			xMove{"go west", "Lobby"},
			xMove{"go east", "Lobby"},
			xMove{"go up", "Parapet"},
			xMove{"go down", "Lobby"}, // first two way
			xMove{"go down", "Basement"},
			xMove{"go up", "Basement"},
			xMove{"go south", "Outside"},
			xMove{"go south", "Foyer"},
			xMove{"enter curtain", "Cloakroom"},
		)
	}
}

type xMove struct {
	cmd, res string
}

func testMoves(t *testing.T, g TestGame, moves ...xMove) {
	// FIX: relations are stored in the model
	if p, ok := g.Model.Instances.FindInstance("player"); assert.True(t, ok) {
		for _, move := range moves {
			t.Logf("%s => %s", move.cmd, move.res)
			out := g.RunInput(move.cmd).FlushOutput()
			if !assert.Equal(t, move.res, where(p)) {
				t.Log(out)
				break
			}
		}
	}
}

func where(inst *M.InstanceInfo) (ret string) {
	if rel, ok := inst.RelativeValue("whereabouts"); ok {
		ret = rel.List()[0]
	}
	return ret
}