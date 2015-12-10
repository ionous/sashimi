package standard

import (
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
)

type ParentLookup struct{ mdl meta.Model }

var objects = ident.MakeId("objects")
var Containment = map[ident.Id]ident.Id{
	"objects-wearer":      "actors-clothing",
	"objects-owner":       "actors-inventory",
	"objects-whereabouts": "rooms-contents",
	"objects-support":     "supporters-contents",
	"objects-enclosure":   "containers-contents",
}

func NewParentLookup(mdl meta.Model) ParentLookup {
	return ParentLookup{mdl}
}
func (p ParentLookup) LookupParent(inst meta.Instance) (ret meta.Instance, rel meta.Property, okay bool) {
	if p.mdl.AreCompatible(inst.GetParentClass().GetId(), objects) {
		for wse, _ := range Containment {
			if prop, ok := inst.GetProperty(wse); ok {
				if parent := prop.GetValue().GetObject(); !parent.Empty() {
					if fini, ok := p.mdl.GetInstance(parent); ok {
						ret, rel, okay = fini, prop, true
						break
					}
				}
			}
		}
	}
	return
}

// FIX: see also: Enclosure. anyway to share?
func (p ParentLookup) LookupRoot(inst meta.Instance) (ret meta.Instance) {
	ret = inst
	for {
		if obj, _, ok := p.LookupParent(ret); !ok {
			break
		} else {
			ret = obj
		}
	}
	return
}
