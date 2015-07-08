package app

import (
	"github.com/ionous/sashimi/util/ident"
)

//
// Record whether a client has heard of a particular object or not.
//
type KnownObjects map[ident.Id]bool

//
// Mark the id'd object as known; return true if newly known.
//
func (this KnownObjects) SetKnown(id ident.Id) (okay bool) {
	if !this[id] {
		this[id] = true
		okay = true
	}
	return okay
}
