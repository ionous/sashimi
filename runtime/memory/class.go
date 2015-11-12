package memory

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/util/ident"
)

type classInfo struct {
	mdl *MemoryModel
	*M.ClassInfo
}

func (c classInfo) GetId() ident.Id {
	return c.Id
}

func (c classInfo) GetParentClass() (ret api.Class) {
	if p := c.Parent; p != nil {
		ret = classInfo{c.mdl, p}
	}
	return
}

func (c classInfo) GetOriginalName() string {
	return c.Plural
}

func (c classInfo) NumProperty() int {
	props := c.mdl.getPropertyList(c.ClassInfo)
	return len(props)
}

func (c classInfo) PropertyNum(i int) api.Property {
	p := c.propertyNum(i)
	return c.getProperty(p)
}

func (c classInfo) propertyNum(i int) M.IProperty {
	props := c.mdl.getPropertyList(c.ClassInfo)
	return props[i] // panics on out of range
}

func (c classInfo) GetProperty(id ident.Id) (ret api.Property, okay bool) {
	// hack for singular and plural properties, note: they wont show up in enumeration...
	var prop M.IProperty
	switch id {
	case plural:
		prop, okay = junkProperty{plural, c.Plural}, true
	case singular:
		prop, okay = junkProperty{singular, c.Singular}, true
	default:
		prop, okay = c.PropertyById(id)
	}
	if okay {
		ret = c.getProperty(prop)
	}
	return
}

func (c classInfo) GetPropertyByChoice(id ident.Id) (ret api.Property, okay bool) {
	if p, _, ok := c.PropertyByChoiceId(id); ok {
		ret, okay = c.getProperty(p), true
	}
	return
}

func (c classInfo) getProperty(p M.IProperty) api.Property {
	return &propBase{
		mdl:      c.mdl,
		src:      c.Id,
		prop:     p,
		getValue: c.getValue,
		setValue: c.setValue}
}

func (c classInfo) getValue(p M.IProperty) GenericValue {
	return p.GetZero(c.Constraints)
}
func (c classInfo) setValue(p M.IProperty, v GenericValue) error {
	panic(fmt.Errorf("classes dont support set property. %s.%v", c.Id, p.GetId()))
}
