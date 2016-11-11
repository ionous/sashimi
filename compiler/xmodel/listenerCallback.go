package xmodel

import (
	"github.com/ionous/mars/rt"
	"github.com/ionous/sashimi/util/ident"
	"github.com/ionous/sashimi/util/sbuf"
	"strconv"
)

// List of event listeners and their callbacks.
type EventCallbacks []EventCallback

type EventCallback struct {
	Event ident.Id
	ListenerCallback
}

// List of action callbacks.
// We cheat by combining action handlers and listener handlers
// there's really not much difference except the phase and,
// probably wont allow actions to change per scene.
type ActionCallbacks []ActionCallback

type ActionCallback struct {
	Action ident.Id
	ListenerCallback
}

// For scenes we dont want to bake these into the class or instances;
// its cheating, but for now using the same structure for class and instance.
// NOTE: Instance can be empty if its a class based listener.
// For the sake of sharing: Even though we listen to events, we point to the action.
type ListenerCallback struct {
	Instance,
	Class ident.Id // Game callback triggered by cb listener.
	Calls   []rt.Execute
	Options ListenerOptions
}

type ListenerOptions int

func (l ListenerOptions) String() string {
	return strconv.FormatInt(int64(l), 16)
}

const (
	EventCapture ListenerOptions = 1 << iota
	EventTargetOnly
	EventQueueAfter
	EventPreventDefault
)

// Create a new class listener: triggers for all instances of the passed class.
func NewClassCallback(
	cls *ClassInfo,
	calls []rt.Execute,
	options ListenerOptions,
) ListenerCallback {
	return ListenerCallback{ident.Empty(), cls.Id, calls, options}
}

// Create a new instance listener: triggers for the passed instance.
func NewInstanceCallback(
	inst *InstanceInfo,
	calls []rt.Execute,
	options ListenerOptions,
) ListenerCallback {
	return ListenerCallback{inst.Id, inst.Class.Id, calls, options}
}

// Return name of instance ( or class ).
func (cb ListenerCallback) GetId() (ret ident.Id) {
	if !cb.Instance.Empty() {
		ret = cb.Instance
	} else {
		ret = cb.Class
	}
	return
}

func (cb ActionCallback) String() string {
	return sbuf.New(cb.GetId(), "->", cb.Action).String()
}

func (cb EventCallback) String() string {
	return sbuf.New(cb.GetId(), "->", cb.Event).String()
}

// UseCapture if the listener wants to participate in the capture cycle ( default is bubble. )
func (cb ListenerCallback) UseCapture() bool {
	return cb.Options&EventCapture != 0
}

// UseTargetOnly if the listener wants callback only when directly targeted.
// ( ie. Event.Target == Event.CurrentTarget )
func (cb ListenerCallback) UseTargetOnly() bool {
	return cb.Options&EventTargetOnly != 0
}

// UseAfterQueue if the listener wants to trigger after a successful event cycle.
func (cb ListenerCallback) UseAfterQueue() bool {
	return cb.Options&EventQueueAfter != 0
}
