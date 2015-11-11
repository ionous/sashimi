package runtime

import (
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/util/ident"
)

type RuntimeConfig struct {
	Calls Callbacks

	// StartFrame and EndFrame should be merged into Output
	// -- and they should be renamed: BeginEvent() EndEvent()
	//*maybe* Target should be mapped into prototype
	// Class should be removed from E.Target
	// only: how do we know that a thing is a "class" and should get "Class" resource?
	// could potentially send target type to startframe
	// right now it seems redicoulous that the game decides that.
	Frame  EventFrame
	Output IOutput
}

type Callbacks interface {
	// LookupCallback returns nil if not found.
	LookupCallback(ident.Id) (G.Callback, bool)
}
