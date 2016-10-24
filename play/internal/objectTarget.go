package internal

import (
	E "github.com/ionous/sashimi/event"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
)

// ObjectTarget implements event.ITarget for Instances.
// The standard rules implement a hierarchy of objects based on containment; for instance: carried object => carrier=> container/supporter of the carrier => room of the contaniner.
type ObjectTarget struct {
	*Dispatch
	obj meta.Instance // FIX? why is the target an instance and not adapter?
}

func NewObjectTarget(d *Dispatch, o meta.Instance) ObjectTarget {
	return ObjectTarget{d, o}
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

// Parent walks up the the (externally defined) containment hierarchy (from event.ITarget.)
func (ot ObjectTarget) Parent() (ret E.ITarget, ok bool) {
	next, haveParent := ot.GetParent(ot.obj)
	cls := ot.obj.GetParentClass()
	if !cls.Empty() || haveParent {
		ret, ok = ClassTarget{ot, cls, next}, true
	}
	return ret, ok
}

// Dispatch an event to an object (from event.ITarget.)
func (ot ObjectTarget) TargetDispatch(evt E.IEvent) (err error) {
	return ot.DispatchEvent(evt, ot.obj.GetId())
}
