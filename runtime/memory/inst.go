package memory

import (
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/util/ident"
)

type instInfo struct {
	mdl *MemoryModel
	*M.InstanceInfo
	class classInfo
}

func (n instInfo) GetId() ident.Id {
	return n.Id
}
func (n instInfo) GetParentClass() api.Class {
	return classInfo{n.mdl, n.Class}
}

func (n instInfo) GetOriginalName() string {
	return n.Name
}

func (n instInfo) NumProperty() int {
	return n.class.NumProperty()
}

func (n instInfo) PropertyNum(i int) (ret api.Property) {
	p := n.class.propertyNum(i)
	return n.getProperty(p)
}

func (n instInfo) GetProperty(id ident.Id) (ret api.Property, okay bool) {
	if p, ok := n.Class.PropertyById(id); ok {
		ret, okay = n.getProperty(p), true
	}
	return
}

func (n instInfo) GetPropertyByChoice(id ident.Id) (ret api.Property, okay bool) {
	if p, _, ok := n.Class.PropertyByChoiceId(id); ok {
		ret, okay = n.getProperty(p), true
	}
	return
}

func (n instInfo) getProperty(p M.IProperty) api.Property {
	return &propBase{
		mdl:      n.mdl,
		src:      n.Id,
		prop:     p,
		getValue: n.getValue,
		setValue: n.setValue}
}

func (n instInfo) getValue(p M.IProperty) (ret GenericValue) {
	// try the object-value interface first
	if v, ok := n.mdl.objectValues.GetValue(n.Id, p.GetId()); ok {
		ret = v
		// fall back to the instance
	} else if v, ok := n.Values[p.GetId()]; ok {
		ret = v
	} else {
		// and from there to class ( chain )
		ret = p.GetZero(n.Class.Constraints)
	}
	return
}

func (n instInfo) setValue(p M.IProperty, v GenericValue) error {
	// STORE FIX: TEST CONSTRAINTS
	return n.mdl.objectValues.SetValue(n.Id, p.GetId(), v)
}
