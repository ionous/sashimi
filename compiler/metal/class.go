package metal

import (
	"fmt"
	M "github.com/ionous/sashimi/compiler/model"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
	"strings"
)

type classInfo struct {
	mdl *Metal
	*M.ClassModel
}

func (c classInfo) GetId() ident.Id {
	return c.Id
}

func (c classInfo) GetParentClass() (ret meta.Class) {
	if p := c.Parent(); !p.Empty() {
		parent := c.mdl.Classes[p]
		ret = classInfo{c.mdl, parent}
	}
	return
}

func (c classInfo) GetOriginalName() string {
	return c.Plural
}

func (c classInfo) NumProperty() int {
	props := c.mdl.getPropertyList(c.ClassModel)
	return len(props)
}

func (c classInfo) PropertyNum(i int) meta.Property {
	p := c.propertyNum(i)
	return c.makeProperty(p)
}

func (c classInfo) propertyNum(i int) *M.PropertyModel {
	props := c.mdl.getPropertyList(c.ClassModel)
	return props[i] // panics on out of range
}

func (c classInfo) GetProperty(id ident.Id) (ret meta.Property, okay bool) {
	if prop, ok := c.getPropertyById(id); ok {
		ret, okay = c.makeProperty(prop), true
	}
	return
}

func (c classInfo) FindProperty(s string) (ret meta.Property, okay bool) {
	if prop, ok := c.getPropertyByName(s); ok {
		ret, okay = c.makeProperty(prop), true
	}
	return
}

func (c classInfo) getPropertyByName(name string) (ret *M.PropertyModel, okay bool) {
	// hack for singular and plural properties, note: they wont show up in enumeration...
	if strings.EqualFold(name, plural.String()) {
		ret, okay = &M.PropertyModel{Id: plural, Type: M.TextProperty}, true
	} else if strings.EqualFold(name, singular.String()) {
		ret, okay = &M.PropertyModel{Id: singular, Type: M.TextProperty}, true
	} else {
		for _, p := range c.mdl.getPropertyList(c.ClassModel) {
			if strings.EqualFold(name, p.Name) {
				ret, okay = p, true
				break
			}
		}
	}
	return
}

func (c classInfo) getPropertyById(id ident.Id) (ret *M.PropertyModel, okay bool) {
	// hack for singular and plural properties, note: they wont show up in enumeration...
	if ident.Compare(id, ident.Join(c.Id, plural)) == 0 {
		ret, okay = &M.PropertyModel{Id: plural, Type: M.TextProperty}, true
	} else if ident.Compare(id, ident.Join(c.Id, singular)) == 0 {
		ret, okay = &M.PropertyModel{Id: singular, Type: M.TextProperty}, true
	} else {
		for _, p := range c.mdl.getPropertyList(c.ClassModel) {
			if id == p.Id {
				ret, okay = p, true
				break
			}
		}
	}
	return
}

func (c classInfo) GetPropertyByChoice(choice ident.Id) (ret meta.Property, okay bool) {
	if p, ok := c.getPropertyByChoice(choice); ok {
		ret, okay = c.makeProperty(p), true
	}
	return
}

func (c classInfo) getPropertyByChoice(id ident.Id) (ret *M.PropertyModel, okay bool) {
	for _, p := range c.mdl.getPropertyList(c.ClassModel) {
		if p.Type == M.EnumProperty {
			if enum, ok := c.mdl.Enumerations[p.Id]; !ok {
				panic(fmt.Sprintf("internal error, couldnt find enumeration '%s.%s'", c.Id, p.Id))
			} else if idx := enum.ChoiceToIndex(id); idx > 0 {
				ret, okay = p, true
				break
			}
		}
	}
	return
}

func (c classInfo) makeProperty(p *M.PropertyModel) meta.Property {
	return &propBase{
		mdl:      c.mdl,
		src:      c.Id,
		prop:     p,
		getValue: c.getValue,
		setValue: c.setValue}
}

func (c classInfo) getValue(p *M.PropertyModel) (ret GenericValue) {
	switch p.Id {
	case plural:
		ret = c.Plural
	case singular:
		ret = c.Singular
	default:
		ret = c.mdl.getZero(p)
	}
	return ret
}

func (c classInfo) setValue(p *M.PropertyModel, v GenericValue) error {
	panic(fmt.Errorf("classes dont support set property. %s.%v", c.Id, p.Id))
}
