package metal

import (
	M "github.com/ionous/sashimi/compiler/model"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
)

// FIX? move to internal so that this can show up in documentation.
type memInst struct {
	mdl *Metal
	*M.InstanceModel
}

func (n *memInst) NumProperty() int {
	return n.getMemClass().NumProperty()
}

func (n *memInst) PropertyNum(i int) (ret meta.Property) {
	p := n.getMemClass().propertyNum(i)
	return n.makeProperty(p)
}

func (n *memInst) GetProperty(id ident.Id) (ret meta.Property, okay bool) {
	if p, ok := n.getMemClass().getPropertyById(id); ok {
		ret, okay = n.makeProperty(p), true
	}
	return
}

func (n *memInst) FindProperty(s string) (ret meta.Property, okay bool) {
	if p, ok := n.getMemClass().getPropertyByName(s); ok {
		ret, okay = n.makeProperty(p), true
	}
	return
}

func (n *memInst) GetPropertyByChoice(id ident.Id) (ret meta.Property, okay bool) {
	if p, ok := n.getMemClass().getPropertyByChoice(id); ok {
		ret, okay = n.makeProperty(p), true
	}
	return
}

func (n *memInst) getMemClass() *memClass {
	cls := n.mdl.Classes[n.Class]
	return &memClass{n.mdl, cls}
}

func (n *memInst) makeProperty(p *M.PropertyModel) meta.Property {
	return makeProperty(n.mdl, p, n)
}

// getStoreId implements valueStore
func (n *memInst) getStoreId() ident.Id {
	return n.Id
}

// getClassId implements valueStore
func (n *memInst) getClassId() ident.Id {
	return n.Class
}

// getValue implements valueStore
func (n *memInst) getValue(slot ident.Id) (ret meta.Generic, okay bool) {
	// try dynamic memory first
	if v, ok := n.mdl.objectValues.GetValue(n.Id, slot); ok {
		ret, okay = v, true
		// fall back to the r/o model
	} else if v, ok := n.Values[slot]; ok {
		ret, okay = v, true
	} else {
		// chain to the class
		ret, okay = n.getMemClass().getValue(slot)
	}
	return
}

// setValue implements valueStore
func (n *memInst) setValue(slot ident.Id, v meta.Generic) (err error) {
	return n.mdl.objectValues.SetValue(n.Id, slot, v)
}
