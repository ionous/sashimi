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

//
// Beacuse of scenes we dont want to bake these into the class or instances;
// its cheating, but for now using the same structure for class and instance.
//
type ListenerCallback struct {
	inst     *InstanceInfo
	class    *ClassInfo
	action   *ActionInfo
	callback Callback
	options  ListenerOptions
}

type ListenerOptions int

const (
	EventCapture ListenerOptions = 1 << iota
	EventTargetOnly
	EventQueueAfter
	EventPreventDefault
)

//
// Create a new class listener: triggers for all instances of the passed class.
//
func NewClassCallback(
	cls *ClassInfo,
	action *ActionInfo,
	callback Callback,
	options ListenerOptions,
) *ListenerCallback {
	return &ListenerCallback{nil, cls, action, callback, options}
}

//
// Create a new instance listener: triggers for the passed instance.
//
func NewInstanceCallback(
	inst *InstanceInfo,
	action *ActionInfo,
	callback Callback,
	options ListenerOptions,
) *ListenerCallback {
	return &ListenerCallback{inst, inst.class, action, callback, options}
}

//
// NOTE: can be nil if its a class based listener.
//
func (this *ListenerCallback) Instance() *InstanceInfo {
	return this.inst
}

//
// Always valid.
//
func (this *ListenerCallback) Class() *ClassInfo {
	return this.class
}

//
// Return name of instance ( or class ).
//
func (this *ListenerCallback) String() string {
	var name string
	if this.inst != nil {
		name = this.inst.name
	} else {
		name = this.class.name
	}
	return fmt.Sprintf("'%s' -> %s", name, this.action)
}

//
// For the sake of sharing: Even though we listen to events, we point to the action.
//
func (this *ListenerCallback) Action() *ActionInfo {
	return this.action
}

//
// Game callback triggered by this listener.
//
func (this *ListenerCallback) Callback() Callback {
	return this.callback
}

// Does this listener want to participate in the capture cycle ( default is bubble. )
//
func (this *ListenerCallback) UseCapture() bool {
	return this.options&EventCapture != 0
}

//
// Does this listener only want callbacks when the event directly targets the listener?
// ( ie. Event.Target == Event.CurrentTarget )
//
func (this *ListenerCallback) UseTargetOnly() bool {
	return this.options&EventTargetOnly != 0
}

//
func (this *ListenerCallback) UseAfterQueue() bool {
	return this.options&EventQueueAfter != 0
}
