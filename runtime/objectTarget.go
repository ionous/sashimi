package runtime

import (
	E "github.com/ionous/sashimi/event"
)

//
// Implements E.ITarget for game objects.
// The standard rules implement a hierarchy of objects based on containment:
// for instance: carried object => carrier=> container/supporter of the carrier => room of the contaniner.
//
type ObjectTarget struct {
	game *Game
	obj  *GameObject // FIX: why is this object and not adapter?
}

//
func (this ObjectTarget) String() string {
	return this.obj.String()
}

//
// Use the ParentLookupStack to walk up the (externally defined) hierarchy.
// (from E.ITarget)
//
func (this ObjectTarget) Parent() (ret E.ITarget, ok bool) {
	game, obj := this.game, this.obj
	cls, next := obj.info.Class(), game.parentLookup.FindParent(obj)
	if cls != nil || next != nil {
		ret, ok = ClassTarget{this, cls, next}, true
	}
	return ret, ok
}

//
// Send an event to this target.
// (from E.ITarget)
//
func (this ObjectTarget) Dispatch(evt E.IEvent) error {
	return this.obj.dispatcher.Dispatch(evt)
}
