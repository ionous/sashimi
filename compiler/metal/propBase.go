package metal

import (
	"fmt"
	M "github.com/ionous/sashimi/compiler/model"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
)

// FIX? to be comparable: would need to stop returning & of propBase,
// and -- i suspect -- would need to change getValue and setValue into non-object methods
type propBase struct {
	mdl  *Metal
	src  ident.Id
	prop *M.PropertyModel
	// life's a little complicated.
	// we have a generic property base ( propBase )
	// an extension to panic on every get and set ( panicValue )
	// and overrides to implement the specific text/num/etc methods ( textValue )
	// the location of values for class and instances differs, so the class and instance pass themselves to their properties, and on to their values.
	getValue func(*M.PropertyModel) GenericValue
	setValue func(*M.PropertyModel, GenericValue) error
}

func (p *propBase) get() GenericValue {
	return p.getValue(p.prop)
}
func (p *propBase) set(v GenericValue) error {
	return p.setValue(p.prop, v)
}

func (p *propBase) String() string {
	return fmt.Sprintf("%s.%s", p.src, p.prop.Id)
}

func (p *propBase) GetId() ident.Id {
	return p.prop.Id
}
func (p *propBase) GetName() string {
	return p.prop.Name
}

func (p *propBase) GetType() meta.PropertyType {
	err := "invalid"
	switch p.prop.Type {
	case M.NumProperty:
		x := meta.NumProperty
		if p.prop.IsMany {
			x |= meta.ArrayProperty
		}
		return x
	case M.TextProperty:
		x := meta.TextProperty
		if p.prop.IsMany {
			x |= meta.ArrayProperty
		}
		return x
	case M.EnumProperty:
		return meta.StateProperty
	case M.PointerProperty:
		x := meta.ObjectProperty
		if p.prop.IsMany {
			x |= meta.ArrayProperty
		}
		return x
	default:
		err = "unknown"
	}
	panic(fmt.Sprintf("GetType(%s.%s) has %s property type %T", p.src, p.prop.Id, err, p.prop))
}

func (p *propBase) GetValue() (ret meta.Value) {
	err := "invalid"
	switch p.prop.Type {
	case M.NumProperty:
		if !p.prop.IsMany {
			return numValue{panicValue{p}}
		}
	case M.TextProperty:
		if !p.prop.IsMany {
			return textValue{panicValue{p}}
		}
	case M.EnumProperty:
		if !p.prop.IsMany {
			return enumValue{panicValue{p}}
		}

	case M.PointerProperty:
		if !p.prop.IsMany {
			// we are not many, since we dont have many to many
			// the far side must be either one or many.
			// many values are "views" onto this objet's own properties
			// one values are real properties, and need to be set.
			if rel, ok := p.mdl.Relations[p.prop.Relation]; ok && rel.Style == M.OneToOne {
				return newRelatedValue(p, rel)
			} else {
				return pointerValue{panicValue{p}}
			}
		}
	default:
		err = "unknown"
	}
	panic(fmt.Sprintf("GetValue(%s.%s) has %s property type %v", p.src, p.prop.Id, err, p.prop.Type))
}

func (p *propBase) GetValues() meta.Values {
	err := "invalid"
	switch p.prop.Type {
	case M.NumProperty:
		if p.prop.IsMany {
			return arrayValues{p, func(i int) meta.Value {
				return numElement{elementValue{panicValue{p}, i}}
			}}
		}
	case M.TextProperty:
		if p.prop.IsMany {
			return arrayValues{p, func(i int) meta.Value {
				return textElement{elementValue{panicValue{p}, i}}
			}}
		}
	case M.EnumProperty:
		//
	case M.PointerProperty:
		if p.prop.IsMany {
			if !p.prop.Relation.Empty() {
				return newManyValues(p)
			} else {
				return arrayValues{p, func(i int) meta.Value {
					return objectElement{elementValue{panicValue{p}, i}}
				}}
			}
		}
	default:
		err = "unknown"
	}
	panic(fmt.Sprintf("GetValues(%s.%s) has %s property type %T", p.src, p.prop.Id, err, p.prop))
}

// FIX: this exists for backwards compatiblity with the client.
// the reality is, a relation effects a table, there may be multiple views that need updating. either the client could do this by seeing the relation and pulling new data,
// or we could push all of thep. this pushes just one. ( client pulling might be best )
func (p *propBase) GetRelative() (ret meta.Relative, okay bool) {
	// get the relation
	if relation, ok := p.mdl.Relations[p.prop.Relation]; ok {
		// get the reverse property
		okay, ret = true, meta.Relative{
			Relation: p.prop.Relation,
			Relates:  p.prop.Relates,
			From:     relation.GetOther(p.prop.Id),
		}
	}
	return
}
