package runtime

import (
	E "github.com/ionous/sashimi/event"
	"github.com/ionous/sashimi/util/ident"
)

// ObjectTarget implements event.ITarget for GameObject.
// The standard rules implement a hierarchy of objects based on containment; for instance: carried object => carrier=> container/supporter of the carrier => room of the contaniner.
type ObjectTarget struct {
	game *Game
	obj  *GameObject // FIX? why is the target object and not adapter?
}

//
func (ot ObjectTarget) Id() ident.Id {
	return ot.obj.Id()
}

//
func (ot ObjectTarget) Class() ident.Id {
	return ot.obj.Class().Id
}

//
func (ot ObjectTarget) String() string {
	return ot.obj.String()
}

// Parent walks up the the (externally defined) containment hierarchy (from event.ITarget.)
func (ot ObjectTarget) Parent() (ret E.ITarget, ok bool) {
	game, obj := ot.game, ot.obj
	cls, next := obj.Class(), game.parentLookup.FindParent(obj)
	if cls != nil || next != nil {
		ret, ok = ClassTarget{ot, cls, next}, true
	}
	return ret, ok
}

// Dispatch an event to an object (from event.ITarget.)
func (ot ObjectTarget) Dispatch(evt E.IEvent) (err error) {
	if d, ok := ot.game.Dispatchers.GetDispatcher(ot.obj.Id()); ok {
		err = d.Dispatch(evt)
	}
	return err
}
