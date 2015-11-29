package app

import (
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
)

// KnownObjects: whether a client has heard of a particular object or not.
type KnownObjects interface {
	SetKnown(meta.Instance) bool
	IsKnown(meta.Instance) bool
}

type KnownObjectMap map[ident.Id]bool

func (known KnownObjectMap) IsKnown(gobj meta.Instance) (okay bool) {
	if gobj != nil {
		id := gobj.GetId()
		okay = known[id]
	}
	return
}

// SetKnown mark the id'd object as known; return true if newly known.
func (known KnownObjectMap) SetKnown(gobj meta.Instance) (okay bool) {
	if gobj != nil {
		if id := gobj.GetId(); !known[id] {
			known[id] = true
			okay = true
		}
	}
	return okay
}
