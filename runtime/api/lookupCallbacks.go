package api

import (
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/util/ident"
)

type LookupCallbacks interface {
	// LookupCallback returns nil if not found.
	LookupCallback(ident.Id) (G.Callback, bool)
}
