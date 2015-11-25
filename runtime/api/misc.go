package api

import (
	E "github.com/ionous/sashimi/event"
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
)

type LookupParents interface {
	// parent instance, property used to find the parent, true if existed
	LookupParent(meta.Model, meta.Instance) (meta.Instance, meta.Property, bool)
}

type LookupCallbacks interface {
	// LookupCallback returns nil if not found.
	LookupCallback(ident.Id) (G.Callback, bool)
}

type EventFrame interface {
	// FUTURE:
	// i just switched to an interface for watching, but i think thats wrong. what game really needs is a way to send events to targets: look at LookupParents whichi is likewise part of this.
	// if the event queue was given to the game as an object perhaps,
	// if it exposed an algorithm.....
	// the other interesting aspect this is how this is attached to LookupCallbacks, since that is how the code calls out.
	//
	//related: *maybe* Target should be mapped into prototype
	// Class should be removed from E.Target;
	// only: how do we know that a thing is a "class" and should get "Class" resource? i had an answer to this one.....
	// could potentially send target type to startframe
	// right now it seems wrong that the game decides that.
	BeginEvent(E.ITarget, E.PathList, *E.Message) IEndEvent
}

type IEndEvent interface {
	RunDefault()
	EndEvent()
}
