package app

import (
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/net/resource"
	"github.com/ionous/sashimi/util/ident"
)

//
func NewObjectRef(gobj meta.Instance) *resource.Object {
	return resource.NewObject(jsonId(gobj.GetId()), jsonId(gobj.GetParentClass()))
}

type RefSerializer struct {
	out   resource.IBuildObjects
	known map[ident.Id]bool
}

func NewRefSerializer(m meta.Model, out resource.ObjectsBuilder) *RefSerializer {
	return &RefSerializer{out, make(map[ident.Id]bool)}
}

func (o *RefSerializer) AddObject(gobj meta.Instance) (ret *resource.Object, okay bool) {
	if id := gobj.GetId(); !o.known[id] {
		o.known[id] = true
		ref := NewObjectRef(gobj)
		o.out.AddObject(ref)
		ret, okay = ref, true
	}
	return
}
