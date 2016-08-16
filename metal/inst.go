package metal

import (
	M "github.com/ionous/sashimi/compiler/model"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
)

// FIX? move to internal so that this can show up in documentation.
// NOTE: because of GameObject equality, cant be pointer receiver.
// now: it can be because of .Equals
type instInfo struct {
	mdl *Metal
	*M.InstanceModel
}

func (n *instInfo) NumProperty() int {
	return n.getClassInfo().NumProperty()
}

func (n *instInfo) PropertyNum(i int) (ret meta.Property) {
	p := n.getClassInfo().propertyNum(i)
	return n.makeProperty(p)
}

func (n *instInfo) GetProperty(id ident.Id) (ret meta.Property, okay bool) {
	if p, ok := n.getClassInfo().getPropertyById(id); ok {
		ret, okay = n.makeProperty(p), true
	}
	return
}

func (n *instInfo) FindProperty(s string) (ret meta.Property, okay bool) {
	if p, ok := n.getClassInfo().getPropertyByName(s); ok {
		ret, okay = n.makeProperty(p), true
	}
	return
}

func (n *instInfo) GetPropertyByChoice(id ident.Id) (ret meta.Property, okay bool) {
	if p, ok := n.getClassInfo().getPropertyByChoice(id); ok {
		ret, okay = n.makeProperty(p), true
	}
	return
}

func (n *instInfo) getClassInfo() *classInfo {
	cls := n.mdl.Classes[n.Class]
	return &classInfo{n.mdl, cls}
}

func (n *instInfo) makeProperty(p *M.PropertyModel) meta.Property {
	return &propBase{
		mdl:      n.mdl,
		src:      n.Id,
		prop:     p,
		getValue: n.getValue,
		setValue: n.setValue}
}

func (n *instInfo) getValue(p *M.PropertyModel) (ret GenericValue) {
	// try the object-value interface first
	if v, ok := n.mdl.objectValues.GetValue(n.Id, p.Id); ok {
		ret = v
		// fall back to the instance
	} else if v, ok := n.Values[p.Id]; ok {
		ret = v
	} else {
		// future: and from there to class defaults ( chain ), currently: zero.
		ret = n.mdl.getZero(p)
	}
	return
}

func (n *instInfo) setValue(p *M.PropertyModel, v GenericValue) error {
	// STORE FIX: TEST CONSTRAINTS
	return n.mdl.objectValues.SetValue(n.Id, p.Id, v)
}
