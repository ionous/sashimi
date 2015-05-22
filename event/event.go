package event

//
// Implemented as interface to help show callers are not intended to modify, but
// since we cant inherit, and i dont want to implement 100 functions for every event
// a generic data() seems required.
//
type IEvent interface {
	Name() string
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
