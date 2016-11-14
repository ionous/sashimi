package api

import (
	"github.com/ionous/sashimi/meta"
)

type LookupParents interface {
	// parent instance, property used to find the parent, true if existed
	LookupParent(meta.Instance) (meta.Instance, error)
}

type ParentHolder struct {
	Parents LookupParents
}

func (p *ParentHolder) LookupParent(i meta.Instance) (ret meta.Instance, err error) {
	if p.Parents != nil {
		ret, err = p.Parents.LookupParent(i)
	}
	return
}
