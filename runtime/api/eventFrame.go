package api

import (
	E "github.com/ionous/sashimi/event"
	"github.com/ionous/sashimi/meta"
)

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
	BeginEvent(meta.Instance, meta.Instance, E.PathList, *E.Message) IEndEvent

	// hacking... there are actions which arent events --
	// this gets the command output to stop consolidating lines across these actions
	FlushFrame()
}

type IEndEvent interface {
	EndEvent()
}
