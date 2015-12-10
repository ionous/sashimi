package app

import (
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
)

type View interface {
	View() ident.Id
	Viewer() ident.Id
	InView(meta.Instance) bool
	ChangedView(meta.Instance, ident.Id, meta.Instance) bool
	EnteredView(meta.Instance, ident.Id, meta.Instance) bool
}
