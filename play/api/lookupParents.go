package api

import (
	"github.com/ionous/sashimi/meta"
)

type LookupParents interface {
	// parent instance, property used to find the parent, true if existed
	LookupParent(meta.Instance) (meta.Instance, meta.Property, bool)
}

type NoParents struct{}

func (NoParents) LookupParent(meta.Instance) (inst meta.Instance, rel meta.Property, okay bool) {
	return
}
