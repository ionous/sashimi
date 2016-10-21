package model

import (
	"github.com/ionous/sashimi/util/ident"
)

// For scenes we dont want to bake these into the class or instances;
// its cheating, but for now using the same structure for class and instance.
// NOTE: Instance can be empty if its a class based listener.
// For the sake of sharing: Even though we listen to events, we point to the action.
type ListenerModel struct {
	Instance,
	Class ident.Id
	Callback CallbackModel // Game callback triggered by cb listener.
	Options  ListenerOptions
}

type ListenerOptions int

const (
	EventTargetOnly ListenerOptions = 1 << iota
	EventQueueAfter
	EventPreventDefault
)

// Return name of instance ( or class ).
func (cb ListenerModel) GetId() (ret ident.Id) {
	if !cb.Instance.Empty() {
		ret = cb.Instance
	} else {
		ret = cb.Class
	}
	return
}

// UseTargetOnly if the listener wants callback only when directly targeted.
// ( ie. Event.Target == Event.CurrentTarget )
func (cb ListenerModel) UseTargetOnly() bool {
	return cb.Options&EventTargetOnly != 0
}

// UseAfterQueue if the listener wants to trigger after a successful event cycle.
func (cb ListenerModel) UseAfterQueue() bool {
	return cb.Options&EventQueueAfter != 0
}
