package script

import "strings"

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
func (this GoingFragment) Through(door string) GoesFromFragment {
	this.fromDoor = door
	return this.GoesFromFragment
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
func (this GoesFromFragment) ConnectsTo(room string) GoesToFragment {
	return GoesToFragment{from: this, toRoom: room, twoWay: true}
}

//
// Establishes a one-way connection between the room From() and the passed destination.
//
func (this GoesFromFragment) ArrivesAt(room string) GoesToFragment {
	return GoesToFragment{from: this, toRoom: room}
}

//
// Optional entrance door in the destination room.
//
func (this GoesToFragment) Door(door string) IFragment {
	this.toDoor = door
	return this
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
func (this GoesToFragment) MakeStatement(b SubjectBlock) (err error) {
	from := newFromSite(b.subject, this.from.fromDoor, this.from.fromDir)
	to := newToSite(this.toRoom, this.toDoor, this.from.fromDir)

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
		} else if this.twoWay {
			_, err = b.The(to.door.str, Has("destination", from.door.str))
		}
	}
	// A Room+Travel Direction (has a matching) Exit
	// ( if you do not have an exit, one will be appointed for you. )
	if err == nil {
		dir := xDir{this.from.fromDir}
		if dir.isSpecified() {
			if _, e := dir.makeDir(b); e != nil {
				err = e
			} else if _, e := b.The(from.room.str, Has(dir.via(), from.door.str)); e != nil {
				err = e
			} else if this.twoWay {
				_, err = b.The(to.room.str, Has(dir.rev(), from.door.str))
			}
		}
	}
	return err
}

// helper to create exit door if needed
func newFromSite(room, door, dir string) xSite {
	if door == "" {
		door = strings.Join([]string{room, "exit", dir}, "-")
	}
	return xSite{xRoom{room}, xDoor{door}}
}

// helper to create entrance door if needed
func newToSite(room, door, dir string) xSite {
	if door == "" {
		door = strings.Join([]string{room, "enter", dir}, "-")
	}
	return xSite{xRoom{room}, xDoor{door}}
}

type xSite struct {
	room xRoom
	door xDoor
}

func (this xSite) makeSite(b SubjectBlock) (err error) {
	if _, e := b.The("room", Called(this.room.str), Exists()); e != nil {
		err = e
	} else {
		_, err = b.The("door", Called(this.door.str), In(this.room.str), Exists())
	}
	return err
}

type xRoom struct {
	str string
}

type xDoor struct {
	str string
}

type xDir struct {
	str string
}

func (this xDir) isSpecified() bool {
	return this.str != ""
}

func (this xDir) makeDir(b SubjectBlock) (int, error) {
	return b.The("direction", Called(this.str), Exists())
}

func (this xDir) via() string {
	return this.str + "-via"
}

func (this xDir) rev() string {
	return this.str + "-rev-via"
}
