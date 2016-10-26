package metal

import (
	M "github.com/ionous/sashimi/compiler/model"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/errutil"
	"github.com/ionous/sashimi/util/ident"
	"strings"
)

type memClass struct {
	mdl *Metal
	*M.ClassModel
}

func (c *memClass) String() string {
	return c.Id.String()
}

func (c *memClass) GetId() ident.Id {
	return c.Id
}

func (c *memClass) GetParentClass() ident.Id {
	if p := c.Parent(); !p.Empty() {
		parent := c.mdl.Classes[p]
		return parent.Id
	}
	return ident.Empty()
}

func (c *memClass) GetOriginalName() string {
	return c.Plural
}

func (c *memClass) NumProperty() int {
	props := c.mdl.propertyList(c.ClassModel)
	return len(props)
}

func (c *memClass) PropertyNum(i int) meta.Property {
	p := c.propertyNum(i)
	return c.makeProperty(p)
}

func (c *memClass) propertyNum(i int) *M.PropertyModel {
	props := c.mdl.propertyList(c.ClassModel)
	return props[i].PropertyModel // panics on out of range
}

func (c *memClass) GetProperty(id ident.Id) (ret meta.Property, okay bool) {
	if prop, ok := c.getPropertyById(id); ok {
		ret, okay = c.makeProperty(prop), true
	}
	return
}

func (c *memClass) FindProperty(s string) (ret meta.Property, okay bool) {
	if prop, ok := c.getPropertyByName(s); ok {
		ret, okay = c.makeProperty(prop), true
	}
	return
}

func (c *memClass) getPropertyByName(name string) (ret *M.PropertyModel, okay bool) {
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

func (c *memClass) getPropertyById(id ident.Id) (ret *M.PropertyModel, okay bool) {
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

func (c *memClass) GetPropertyByChoice(choice ident.Id) (ret meta.Property, okay bool) {
	if p, ok := c.getPropertyByChoice(choice); ok {
		ret, okay = c.makeProperty(p), true
	}
	return
}

func (c *memClass) getPropertyByChoice(id ident.Id) (ret *M.PropertyModel, okay bool) {
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

func (c *memClass) makeProperty(p *M.PropertyModel) meta.Property {
	return makeProperty(c.mdl, p, c)
}

// getStoreId implements valueStore
func (c *memClass) getStoreId() ident.Id {
	return ident.Empty()
}

// getClassId implements valueStore
func (c *memClass) getClassId() ident.Id {
	return c.Id
}

// getValue implements valueStore
func (c *memClass) getValue(slot ident.Id) (ret meta.Generic, okay bool) {
	switch slot {
	case pluralId:
		ret, okay = c.Plural, true
	case singularId:
		ret, okay = c.Singular, true
	default:
		// MARS: implement class defaults
	}
	return
}

// setValue implements valueStore
func (c *memClass) setValue(slot ident.Id, v meta.Generic) error {
	return errutil.New("classes dont support set property")
}
