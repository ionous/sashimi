package event

import "github.com/ionous/sashimi/util/ident"

// IEvent
// implemented as interface to help show callers are not intended to modify, but
// since we cant inherit, and i dont want to implement 100 functions for every event
// a generic data() seems required.
type IEvent interface {
	Id() ident.Id
	Data() interface{}

	Bubbles() bool
	Cancelable() bool
	DefaultBlocked() bool

	Target() ITarget // ultimate endpoint for the event
	Path() PathList

	Phase() Phase
	CurrentTarget() ITarget // which changes as the event captures and bubbles

	// cancel the default action
	// returns true if now cancelled
	PreventDefault() bool
	// stop processing the event flow after this set is done
	StopPropagation()
	// stop processing all other event handlers immediately
	StopImmediatePropagation()
}

// IListen in order to handle event callbacks.
// Uses an interface for facilitating add/remove event listeners;
// comparing function pointers is error-prone in go (due to closures)
type IListen interface {
	// FIX: does dispatch really need an error handling?
	HandleEvent(IEvent) error
}

// ITarget dispatch events to some hierarchical node, for instance, in a DOM.
// FIX: i wonder whether instead of having ITarget ( target could be an id )
// the focus could be around Path: rather than building ( and copying ) the path ahead of time, instead: have a path interface which the caller implements.
// dispatch could be to a path node.
// evaluate what dispatch needs, not what we want to give it.
type ITarget interface {
	Id() ident.Id
	Parent() (ITarget, bool)     // used by path.addPath
	TargetDispatch(IEvent) error // used by proc.sendToTarget, Send(path PathList)
}
