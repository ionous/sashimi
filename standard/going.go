package script

import (
	"strings"
)

var Directions = []string{"north", "south", "east", "west", "up", "down"}

// FIX: move these into a standard rules extension package?
func _makeOpposites() map[string]string {
	op := make(map[string]string)
	for i := 0; i < len(Directions)/2; i++ {
		a, b := Directions[2*i], Directions[2*i+1]
		op[a], op[b] = b, a
	}
	return op
}

var opposites = _makeOpposites()

//
// Begin a statement connecting one room to another via a movement direction.
// Direction: a noun of "direction" type: ex. north, east, south.
//
func Going(direction string) GoingFragment {
	return GoingFragment{GoesFromFragment{origin: NewOrigin(2), fromDir: direction}}
}

//
// Causes directional movement to pass through an explicit departure door.
//
func (goingFrom GoingFragment) Through(door string) GoesFromFragment {
	goingFrom.fromDoor = door
	return goingFrom.GoesFromFragment
}

//
// Begin a statement connecting one room to another via a door.
// Door: The exit from a room.
//
func Through(door string) GoesFromFragment {
	return GoesFromFragment{origin: NewOrigin(2), fromDoor: door}
}

//
// Establishes a two-way connection between the room From() and the passed destination.
//
func (goesFrom GoesFromFragment) ConnectsTo(room string) GoesToFragment {
	return GoesToFragment{from: goesFrom, toRoom: room, twoWay: true}
}

//
// Establishes a one-way connection between the room From() and the passed destination.
//
func (goesFrom GoesFromFragment) ArrivesAt(room string) GoesToFragment {
	return GoesToFragment{from: goesFrom, toRoom: room}
}

//
// Optional door to arrive at in the destination room.
//
func (goesTo GoesToFragment) Door(door string) IFragment {
	goesTo.toDoor = door
	return goesTo
}

type GoingFragment struct {
	GoesFromFragment
}

type GoesFromFragment struct {
	origin            Origin
	fromDir, fromDoor string
}

type GoesToFragment struct {
	from           GoesFromFragment
	toRoom, toDoor string
	twoWay         bool
}

//
// implements IFragment for use in The()
//
func (goesTo GoesToFragment) MakeStatement(b SubjectBlock) (err error) {
	from := newFromSite(b.subject, goesTo.from.fromDoor, goesTo.from.fromDir)
	to := newToSite(goesTo.toRoom, goesTo.toDoor, goesTo.from.fromDir)

	// A Room (contains) Doors
	if e := from.makeSite(b); e != nil {
		err = e
	} else if e := to.makeSite(b); e != nil {
		err = e
	}
	// A departure door (has a matching) arrival door
	if err == nil {
		if _, e := b.The(from.door.str, Has("destination", to.door.str)); e != nil {
			err = e
		} else if goesTo.twoWay {
			_, err = b.The(to.door.str, Has("destination", from.door.str))
		}
	}
	// A Room+Travel Direction (has a matching) departure door
	// ( if you do not have an deptature door, one will be created for you. )
	if err == nil {
		dir := xDir{goesTo.from.fromDir}
		if dir.isSpecified() {
			if _, e := dir.makeDir(b); e != nil {
				err = e
			} else if _, e := b.The(from.room.str, Has(dir.via(), from.door.str)); e != nil {
				err = e
			} else if goesTo.twoWay {
				_, err = b.The(to.room.str, Has(dir.revVia(), to.door.str))
				// FIX? REMOVED dynamic opposite lookup
				// needs more thought as to how new directions could be added
				// perhaps some sort of "dependency injection" where we can add evaluations
				// -- dynmic compiler generators -- as hooks after ( dependent on ) sets of other instances, classes, etc. so those hooks can use model reflection to generate new, non-conflicting, model data -- this is already similar to the idea of onion skins of visual content, hardpoint hooks, etc.
				//_, err = b.The(to.room.str, Has(dir.revRev(), from.door.str))
			}
		}
	}
	return err
}

// helper to create departure door if needed
func newFromSite(room, door, dir string) xSite {
	gen := door == ""
	if gen {
		door = strings.Join([]string{room, "departure", dir}, "-")
	}
	return xSite{xRoom{room}, xDoor{door, gen}}
}

// helper to create arrival door if needed
func newToSite(room, door, dir string) xSite {
	gen := door == ""
	if gen {
		door = strings.Join([]string{room, "arrival", dir}, "-")
	}
	return xSite{xRoom{room}, xDoor{door, gen}}
}

type xSite struct {
	room xRoom
	door xDoor
}

func (x xSite) makeSite(b SubjectBlock) (err error) {
	if _, e := b.The("room", Called(x.room.str), Exists()); e != nil {
		err = e
	} else if _, e = b.The("door", Called(x.door.str), In(x.room.str), Exists()); e != nil {
		err = e
	} else if x.door.gen {
		_, err = b.Our(x.door.str, Is("scenery"))
	}
	return err
}

type xRoom struct {
	str string
}

type xDoor struct {
	str string
	gen bool
}

type xDir struct {
	str string
}

func (x xDir) isSpecified() bool {
	return len(x.str) > 0
}

func (x xDir) makeDir(b SubjectBlock) (int, error) {
	return b.The("direction", Called(x.str), Exists())
}

func (x xDir) via() string {
	return x.str + "-via"
}
func (x xDir) opposite() string {
	return opposites[x.str]
}

func (x xDir) revVia() string {
	return x.opposite() + "-via"
}

// FIX? REMOVED dynamic opposite lookup ( see comment in MakeStatement )
// func (x xDir) revVia() string {
// 	return x.str + "-rev-via"
// }
