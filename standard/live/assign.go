package live

import G "github.com/ionous/sashimi/game"

// FIX: there's no error testing here and its definitely possible to screw things up.
func AssignTo(src G.IObject, rel string, dest G.IObject) {
	// sure hope there's no errors, would relation by value remove the need for transaction?
	if _, parentRel := src.ParentRelation(); len(parentRel) > 0 {
		// note: objects which start out of world, dont have an owner field to clear.
		src.Set(parentRel, nil)
	}
	src.Get(rel).SetObject(dest)
}
