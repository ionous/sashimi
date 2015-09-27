package event

import "github.com/ionous/sashimi/util/ident"

//
// Event handler callback.
// Uses an interface for facilitating add/remove event listeners;
// comparing function pointers is error-prone in go (due to closures)
//
type IListen interface {
	HandleEvent(IEvent) error
}

//
// Node, for instance, in a DOM.
//
type ITarget interface {
	Id() ident.Id
	Class() ident.Id
	Parent() (ITarget, bool)
	Dispatch(IEvent) error
}
