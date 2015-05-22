package event

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
	Parent() (ITarget, bool)
	Dispatch(IEvent) error
}
