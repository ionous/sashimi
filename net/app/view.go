package app

import (
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
)

type View interface {
	Viewpoint() meta.Instance
	InView(meta.Instance) bool
	ChangedView(meta.Instance, ident.Id, meta.Instance) bool
	EnteredView(meta.Instance, ident.Id, meta.Instance) bool
}
