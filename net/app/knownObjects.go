package app

import (
	M "github.com/ionous/sashimi/model"
)

//
// Record whether a client has heard of a particular object or not.
//
type KnownObjects map[M.StringId]bool

//
// Mark the id'd object as known; return true if newly known.
//
func (this KnownObjects) SetKnown(id M.StringId) (okay bool) {
	if !this[id] {
		this[id] = true
		okay = true
	}
	return okay
}
