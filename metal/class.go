package metal

import (
	M "github.com/ionous/sashimi/compiler/model"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
	"strings"
)

type classInfo struct {
	mdl *Metal
	*M.ClassModel
}

func (c *classInfo) String() string {
	return c.Id.String()
}

func (c *classInfo) GetId() ident.Id {
	return c.Id
}

func (c *classInfo) GetParentClass() ident.Id {
	if p := c.Parent(); !p.Empty() {
		parent := c.mdl.Classes[p]
		return parent.Id
	}
	return ident.Empty()
}

func (c *classInfo) GetOriginalName() string {
	return c.Plural
}

func (c *classInfo) NumProperty() int {
	props := c.mdl.propertyList(c.ClassModel)
	return len(props)
}

func (c *classInfo) PropertyNum(i int) meta.Property {
	p := c.propertyNum(i)
	return c.makeProperty(p)
}

func (c *classInfo) propertyNum(i int) *M.PropertyModel {
	props := c.mdl.propertyList(c.ClassModel)
	return props[i].PropertyModel // panics on out of range
}

func (c *classInfo) GetProperty(id ident.Id) (ret meta.Property, okay bool) {
	if prop, ok := c.getPropertyById(id); ok {
		ret, okay = c.makeProperty(prop), true
	}
	return
}

func (c *classInfo) FindProperty(s string) (ret meta.Property, okay bool) {
	if prop, ok := c.getPropertyByName(s); ok {
		ret, okay = c.makeProperty(prop), true
	}
	return
}

func (c *classInfo) getPropertyByName(name string) (ret *M.PropertyModel, okay bool) {
	// FIX: hack for singular and plural properties, note: they wont show up in enumeration...
	// ie. asking for FindProperty("plural"), or FindProperty("singular")
	// these really should be generated at compile time or something
	lower := strings.ToLower(name)
	if lower == pluralString {
		ret, okay = &M.PropertyModel{Id: pluralId, Type: M.TextProperty}, true
	} else if lower == singularString {
		ret, okay = &M.PropertyModel{Id: singularId, Type: M.TextProperty}, true
	} else {
		for _, p := range c.mdl.propertyList(c.ClassModel) {
			if lower == p.lower {
				ret, okay = p.PropertyModel, true
				break
			}
		}
	}

	return
}

func (c *classInfo) getPropertyById(id ident.Id) (ret *M.PropertyModel, okay bool) {
	// hack for singular and plural properties, note: they wont show up in enumeration...
	if id.Equals(ident.Join(c.Id, pluralId)) {
		ret, okay = &M.PropertyModel{Id: pluralId, Type: M.TextProperty}, true
	} else if id.Equals(ident.Join(c.Id, singularId)) {
		ret, okay = &M.PropertyModel{Id: singularId, Type: M.TextProperty}, true
	} else {
		for _, p := range c.mdl.propertyList(c.ClassModel) {
			if id == p.Id {
				ret, okay = p.PropertyModel, true
				break
			}
		}
	}
	return
}

func (c *classInfo) GetPropertyByChoice(choice ident.Id) (ret meta.Property, okay bool) {
	if p, ok := c.getPropertyByChoice(choice); ok {
		ret, okay = c.makeProperty(p), true
	}
	return
}

func (c *classInfo) getPropertyByChoice(id ident.Id) (ret *M.PropertyModel, okay bool) {
	for _, p := range c.mdl.propertyList(c.ClassModel) {
		if p.Type == M.EnumProperty {
			if enum, ok := c.mdl.Enumerations[p.Id]; ok {
				if idx := enum.ChoiceToIndex(id); idx > 0 {
					ret, okay = p.PropertyModel, true
					break
				}
			}
		}
	}
	return
}

func (c *classInfo) makeProperty(p *M.PropertyModel) meta.Property {
	return &propBase{
		mdl:      c.mdl,
		src:      c.Id,
		prop:     p,
		getValue: c.getValue,
		setValue: c.setValue}
}

func (c *classInfo) getValue(p *M.PropertyModel) (ret meta.Generic) {
	switch p.Id {
	case pluralId:
		ret = c.Plural
	case singularId:
		ret = c.Singular
	default:
		ret = c.mdl.getZero(p)
	}
	return
}

func (c *classInfo) setValue(p *M.PropertyModel, v meta.Generic) error {
	// test current expect full on panic, even through we return an error... hmmm...
	panic("classes dont support set property")
}
