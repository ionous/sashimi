package standard

import (
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
)

type ParentLookup struct{}

var objects = ident.MakeId("objects")
var containment = []string{"wearer", "owner", "whereabouts", "support", "enclosure"}

func (p ParentLookup) LookupParent(mdl meta.Model, inst meta.Instance) (ret meta.Instance, rel meta.Property, okay bool) {
	if mdl.AreCompatible(inst.GetParentClass().GetId(), objects) {
		for _, wse := range containment {
			if prop, ok := inst.FindProperty(wse); ok {
				if parent := prop.GetValue().GetObject(); !parent.Empty() {
					if fini, ok := mdl.GetInstance(parent); ok {
						ret, rel, okay = fini, prop, true
						break
					}
				}
			}
		}
	}
	return
}
