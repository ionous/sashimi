package api

import (
	"github.com/ionous/sashimi/meta"
)

type LookupParents interface {
	// parent instance, property used to find the parent, true if existed
	LookupParent(meta.Instance) (meta.Instance, meta.Property, bool)
}
