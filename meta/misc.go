package api

import (
	E "github.com/ionous/sashimi/event"
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/util/ident"
)

// NewEventFrame returns a function for defer() end of event.
// FIX? does the event data need to be copied as well?
// FIX: this is bit uglys gross.
type EventFrame func(E.ITarget, *E.Message) func()

type LookupParents interface {
	// parent instance, property used to find the parent, true if existed
	LookupParent(Model, Instance) (Instance, Property, bool)
}

type LookupCallbacks interface {
	// LookupCallback returns nil if not found.
	LookupCallback(ident.Id) (G.Callback, bool)
}
