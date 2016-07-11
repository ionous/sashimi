package live

import (
	G "github.com/ionous/sashimi/game"
)

// FIX: there's no error testing here and its definitely possible to screw things up.
func AssignTo(prop G.IObject, rel string, dest G.IObject) {
	// sure hope there's no errors, would relation by value remove the need for transaction?
	if _, parentRel := prop.ParentRelation(); parentRel != "" {
		// note: an object like the fishFood isnt "in the world", and doest have an owner field to clear.
		prop.Set(parentRel, nil)
	}
	prop.Set(rel, dest)
}
