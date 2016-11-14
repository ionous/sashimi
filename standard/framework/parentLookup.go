package framework

import (
	"github.com/ionous/mars/rt"
	"github.com/ionous/mars/std"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/errutil"
)

type ParentLookup struct {
	Run rt.Runtime
}

func NewParentLookup() *ParentLookup {
	return &ParentLookup{}
}

func (p *ParentLookup) LookupParent(i meta.Instance) (ret meta.Instance, err error) {
	if p.Run == nil {
		err = errutil.New("parent lookup has no runtime")
	} else {
		if i, e := std.Parent(rt.Object{i}).GetObject(p.Run); e != nil {
			err = e
		} else {
			ret = i.Instance
		}
	}
	return
}
