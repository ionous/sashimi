package memory

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/util/ident"
)

type propBase struct {
	mdl  *MemoryModel
	src  ident.Id
	prop M.IProperty
	// life's a little complicated.
	// we have a generic property base ( propBase )
	// an extension to panic on every get and set ( panicValue )
	// and overrides to implement the specific text/num/etc methods ( textValue )
	// the location of values for class and instances differs, so the class and instance pass themselves to their properties, and on to their values.
	getValue func(M.IProperty) GenericValue
	setValue func(M.IProperty, GenericValue) error
}

func (p *propBase) get() GenericValue {
	return p.getValue(p.prop)
}
func (p *propBase) set(v GenericValue) error {
	return p.setValue(p.prop, v)
}

func (p *propBase) String() string {
	return fmt.Sprintf("%s.%s", p.src, p.prop.GetId())
}

func (p *propBase) GetId() ident.Id {
	return p.prop.GetId()
}

func (p *propBase) GetType() api.PropertyType {
	err := "invalid"
	switch m := p.prop.(type) {
	case M.NumProperty:
		x := api.NumProperty
		if m.IsMany {
			x |= api.ArrayProperty
		}
		return x
	case M.TextProperty:
		x := api.TextProperty
		if m.IsMany {
			x |= api.ArrayProperty
		}
		return x
	case junkProperty:
		x := api.TextProperty
		return x
	case M.EnumProperty:
		return api.StateProperty
	case M.PointerProperty:
		x := api.ObjectProperty
		if m.IsMany {
			x |= api.ArrayProperty
		}
		return x
	case M.RelativeProperty:
		x := api.ObjectProperty
		if m.IsMany {
			x |= api.ArrayProperty
		}
		return x
	default:
		err = "unknown"
	}
	panic(fmt.Sprintf("GetType(%s.%s) has %s property type %T", p.src, p.prop.GetId(), err, p.prop))
}

func (p *propBase) GetValue() api.Value {
	err := "invalid"
	switch m := p.prop.(type) {
	case M.NumProperty:
		if !m.IsMany {
			return numValue{panicValue{p}}
		}
	case M.TextProperty:
		if !m.IsMany {
			return textValue{panicValue{p}}
		}
	case junkProperty:
		return textValue{panicValue{p}}

	case M.EnumProperty:
		return enumValue{panicValue{p}}

	case M.PointerProperty:
		if !m.IsMany {
			return pointerValue{panicValue{p}}
		}
	case M.RelativeProperty:
		if !m.IsMany {
			return singleValue(p)
		}
	default:
		err = "unknown"
	}
	panic(fmt.Sprintf("GetValue(%s.%s) has %s property type %T", p.src, p.prop.GetId(), err, p.prop))
}

func (p *propBase) GetValues() api.Values {
	err := "invalid"
	switch m := p.prop.(type) {
	case M.NumProperty:
		if m.IsMany {
			return arrayValues{p, func(i int) api.Value {
				return numElement{elementValue{panicValue{p}, i}}
			}}
		}
	case M.TextProperty:
		if m.IsMany {
			return arrayValues{p, func(i int) api.Value {
				return textElement{elementValue{panicValue{p}, i}}
			}}
		}
	case M.EnumProperty:
		//
	case M.PointerProperty:
		if m.IsMany {
			return arrayValues{p, func(i int) api.Value {
				return objectElement{elementValue{panicValue{p}, i}}
			}}
		}
	case M.RelativeProperty:
		if m.IsMany {
			return manyValue(p)
		}
	default:
		err = "unknown"
	}
	panic(fmt.Sprintf("GetValues(%s.%s) has %s property type %T", p.src, p.prop.GetId(), err, p.prop))
}

func (p *propBase) GetRelative() (ret api.Relative, okay bool) {
	switch prop := p.prop.(type) {
	case M.PointerProperty:
	case M.RelativeProperty:

		// get the relation
		relation := p.mdl.Relations[prop.Relation]

		// get the reverse property
		other := relation.GetOther(prop.IsRev)

		okay, ret = true, api.Relative{
			Relation: prop.Relation,
			Relates:  prop.Relates,
			// FIX: this exists for backwards compatiblity with the client.
			// the reality is, a relation effects a table, there may be multiple views that need updating. either the client could do this by seeing the relation and pulling new data,
			// or we could push all of them. this pushes just one. ( client pulling might be best )
			From:  other.Property,
			IsRev: prop.IsRev,
		}
	}
	return
}
