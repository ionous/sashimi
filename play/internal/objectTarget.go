package internal

import (
	"github.com/ionous/mars/rt"
	"github.com/ionous/mars/rtm"
	E "github.com/ionous/sashimi/event"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
)

// ObjectTarget implements event.ITarget for Instances.
// The standard rules implement a hierarchy of objects based on containment; for instance: carried object => carrier=> container/supporter of the carrier => room of the contaniner.
type ObjectTarget struct {
	*rtm.ActionRuntime
	obj meta.Instance // FIX? why is the target an instance and not adapter?
}

func NewObjectTarget(act *rtm.ActionRuntime, o meta.Instance) ObjectTarget {
	return ObjectTarget{act, o}
}

//
func (ot ObjectTarget) Id() ident.Id {
	return ot.obj.GetId()
}

//
func (ot ObjectTarget) Class() ident.Id {
	return ot.obj.GetParentClass()
}

//
func (ot ObjectTarget) String() string {
	return ot.obj.GetId().String()
}

// Parent walks up the the (externally defined) containment hierarchy first exhausting the classes of this object. (see also: event.ITarget.)
func (ot ObjectTarget) Parent() (ret E.ITarget, ok bool) {
	// MARS: handle error
	next, err := ot.FindParent(rt.Object{ot.obj})
	if err != nil {
		panic(err)
	}
	cls := ot.obj.GetParentClass()
	if !cls.Empty() || (next.Instance != nil) {
		ret, ok = ClassTarget{ot, cls, next.Instance}, true
	}
	return ret, ok
}

// Dispatch an event to an object (from event.ITarget.)
func (ot ObjectTarget) TargetDispatch(evt E.IEvent) (err error) {
	return ot.DispatchEvent(evt, ot.obj.GetId())
}

// note: we get multiple dispatch calls for each event: capture, target, and bubble.
func (ot ObjectTarget) DispatchEvent(evt E.IEvent, target ident.Id) (err error) {
	if src, ok := ot.GetEvent(evt.Id()); ok {
		if ls, ok := src.GetListeners(true); ok {
			err = E.Capture(evt, NewGameListeners(ot.ActionRuntime, evt, target, ls))
		}
		if err == nil {
			if ls, ok := src.GetListeners(false); ok {
				err = E.Bubble(evt, NewGameListeners(ot.ActionRuntime, evt, target, ls))
			}
		}
	}
	return
}
