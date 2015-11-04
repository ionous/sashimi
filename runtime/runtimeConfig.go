package runtime

import (
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/util/ident"
)

type RuntimeConfig struct {
	Calls  Callbacks
	Frame  EventFrame
	Output IOutput
}

type Callbacks interface {
	// Lookup returns nil if not found.
	Lookup(ident.Id) G.Callback
}
