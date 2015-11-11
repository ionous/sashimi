package runtime

import (
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/util/ident"
)

type PropertyChange interface {
	NumChange(obj api.Instance, prop ident.Id, prev, next float32)
	TextChange(obj api.Instance, prop ident.Id, prev, next string)
	StateChange(obj api.Instance, prop ident.Id, prev, next ident.Id)
	ReferenceChange(obj api.Instance, prop, other ident.Id, prev, next api.Instance)
}
