package call

import (
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/util/ident"
)

type Compiler interface {
	CompileCallback(G.Callback) (ident.Id, error)
}
