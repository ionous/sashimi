package runtime

import (
	G "github.com/ionous/sashimi/game"
	M "github.com/ionous/sashimi/model"
)

type Config struct {
	Calls  Callbacks
	Frame  EventFrame
	Output IOutput
}

type Callbacks interface {
	// Lookup returns nil if not found.
	Lookup(M.Callback) G.Callback
}
