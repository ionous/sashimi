package app

import (
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/util/ident"
)

//
// Record whether a client has heard of a particular object or not.
//
type KnownObjects map[ident.Id]bool

//
// Mark the id'd object as known; return true if newly known.
//
func (known KnownObjects) SetKnown(gobj api.Instance) (okay bool) {
	if gobj != nil {
		if id := gobj.GetId(); !known[id] {
			known[id] = true
			okay = true
		}
	}
	return okay
}
