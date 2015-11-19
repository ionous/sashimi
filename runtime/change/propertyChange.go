package runtime

import (
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
)

type PropertyChange interface {
	NumChange(obj meta.Instance, prop ident.Id, prev, next float32)
	TextChange(obj meta.Instance, prop ident.Id, prev, next string)
	StateChange(obj meta.Instance, prop ident.Id, prev, next ident.Id)
	ReferenceChange(obj meta.Instance, prop, other ident.Id, prev, next meta.Instance)
}
