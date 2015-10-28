package model

import (
	"fmt"
)

// List of event listeners and their callbacks.
type ListenerCallbacks []*ListenerCallback

// List of action callbacks.
// We cheat by combining action handlers and listener handlers
// there's really not much difference except the phase and,
// probably wont allow actions to change per scene.
type ActionCallbacks []*ListenerCallback

// Beacuse of scenes we dont want to bake these into the class or instances;
// its cheating, but for now using the same structure for class and instance.
type ListenerCallback struct {
	Instance *InstanceInfo // NOTE: can be nil if its a class based listener.
	Class    *ClassInfo    // Always valid.
	Action   *ActionInfo   // For the sake of sharing: Even though we listen to events, we point to the action.

	Callback Callback // Game callback triggered by this listener.
	Options  ListenerOptions
}

type ListenerOptions int

const (
	EventCapture ListenerOptions = 1 << iota
	EventTargetOnly
	EventQueueAfter
	EventPreventDefault
)

// Create a new class listener: triggers for all instances of the passed class.
func NewClassCallback(
	cls *ClassInfo,
	action *ActionInfo,
	callback Callback,
	options ListenerOptions,
) *ListenerCallback {
	return &ListenerCallback{nil, cls, action, callback, options}
}

// Create a new instance listener: triggers for the passed instance.
func NewInstanceCallback(
	inst *InstanceInfo,
	action *ActionInfo,
	callback Callback,
	options ListenerOptions,
) *ListenerCallback {
	return &ListenerCallback{inst, inst.Class, action, callback, options}
}

// Return name of instance ( or class ).
func (this *ListenerCallback) String() string {
	var name string
	if this.Instance != nil {
		name = this.Instance.Name
	} else {
		name = this.Class.Plural
	}
	return fmt.Sprintf("'%s' -> %s", name, this.Action)
}

// Does this listener want to participate in the capture cycle ( default is bubble. )
func (this *ListenerCallback) UseCapture() bool {
	return this.Options&EventCapture != 0
}

// Does this listener only want callbacks when the event directly targets the listener?
// ( ie. Event.Target == Event.CurrentTarget )
func (this *ListenerCallback) UseTargetOnly() bool {
	return this.Options&EventTargetOnly != 0
}

//
func (this *ListenerCallback) UseAfterQueue() bool {
	return this.Options&EventQueueAfter != 0
}
