package call

import (
	G "github.com/ionous/sashimi/game"
)

type Compiler interface {
	Compile(G.Callback) (Marker, error)
}
