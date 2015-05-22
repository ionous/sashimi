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
	cls, next := obj.Info().Class(), game.parentLookup.FindParent(obj)
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

// func (this ObjectTarget) Dispatch(evt E.IEvent) (err error) {
// 	// note: if something isnt firing; double check the name.
// 	// make sure it's not an action name, but the event-ing name.
// 	game, obj := this.game, this.obj
// 	phase := evt.Phase()

// 	// capturing or targeting? trigger class capture listeners
// 	if phase != E.BubblingPhase {
// 		err = game.Dispatchers.DispatchClassEvent(evt, obj.info.Class(), true)
// 	}

// 	if err == nil {
// 		// instance listeners:
// 		if e := this.dispatcher.Dispatch(evt); e != nil {
// 			err = e
// 		} else {
// 			// bubbling or targeting? trigger class bubble listeners
// 			if err == nil && phase != CapturingPhase {
// 				err = game.Dispatchers.DispatchClassEvent(evt, obj.info.Class(), false)
// 			}
// 		}
// 	}
// 	return err
// }
