package appengine

import (
	"github.com/ionous/sashimi/meta"
)

// net/app/KnownObjects
type AlwaysKnown struct{}

// IsKnown always returns true.
func (k AlwaysKnown) IsKnown(gobj meta.Instance) bool {
	return gobj != nil
}

// SetKnown always returns true: as if the object was always knewly known.
func (k AlwaysKnown) SetKnown(gobj meta.Instance) bool {
	return gobj != nil
}
