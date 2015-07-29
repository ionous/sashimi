package script

import (
	"strings"
)

// FIX: move these into a standard rules extension package?

//
// Begin a statement connecting one room to another via a movement direction.
// Direction: a noun of "direction" type: ex. north, east, south.
//
func Going(direction string) GoingFragment {
	return GoingFragment{GoesFromFragment{origin: NewOrigin(2), fromDir: direction}}
}

//
// Causes directional movement to pass through an explicit exit door.
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
// Optional entrance door in the destination room.
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
	// An Exit (has a matching) Entrance
	if err == nil {
		if _, e := b.The(from.door.str, Has("destination", to.door.str)); e != nil {
			err = e
		} else if goesTo.twoWay {
			_, err = b.The(to.door.str, Has("destination", from.door.str))
		}
	}
	// A Room+Travel Direction (has a matching) Exit
	// ( if you do not have an exit, one will be appointed for you. )
	if err == nil {
		dir := xDir{goesTo.from.fromDir}
		if dir.isSpecified() {
			if _, e := dir.makeDir(b); e != nil {
				err = e
			} else if _, e := b.The(from.room.str, Has(dir.via(), from.door.str)); e != nil {
				err = e
			} else if goesTo.twoWay {
				_, err = b.The(to.room.str, Has(dir.rev(), from.door.str))
			}
		}
	}
	return err
}

// helper to create exit door if needed
func newFromSite(room, door, dir string) xSite {
	gen := door == ""
	if gen {
		door = strings.Join([]string{room, "exit", dir}, "-")
	}
	return xSite{xRoom{room}, xDoor{door, gen}}
}

// helper to create entrance door if needed
func newToSite(room, door, dir string) xSite {
	gen := door == ""
	if gen {
		door = strings.Join([]string{room, "enter", dir}, "-")
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
	return x.str != ""
}

func (x xDir) makeDir(b SubjectBlock) (int, error) {
	return b.The("direction", Called(x.str), Exists())
}

func (x xDir) via() string {
	return x.str + "-via"
}

func (x xDir) rev() string {
	return x.str + "-rev-via"
}
