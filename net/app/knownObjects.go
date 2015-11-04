package app

import (
	R "github.com/ionous/sashimi/runtime"
	"github.com/ionous/sashimi/util/ident"
)

//
// Record whether a client has heard of a particular object or not.
//
type KnownObjects map[ident.Id]bool

//
// Mark the id'd object as known; return true if newly known.
//
func (known KnownObjects) SetKnown(gobj *R.GameObject) (okay bool) {
	if gobj != nil {
		if id := gobj.Id(); !known[id] {
			known[id] = true
			okay = true
		}
	}
	return okay
}
