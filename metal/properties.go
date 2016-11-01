package metal

import (
	"github.com/ionous/mars/rt"
	M "github.com/ionous/sashimi/compiler/model"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/errutil"
	"github.com/ionous/sashimi/util/ident"
	"github.com/ionous/sashimi/util/sbuf"
)

// two implentations: memClass and memInst
type valueStore interface {
	// relation compatibility tests need to know parent class
	getClassId() ident.Id
	// relation reflexive storage ( storing this in that ) needs instance id.
	getStoreId() ident.Id
	// caller can verify property is of expected type if desired;
	// we dont need to do that as well.
	getValue(slot ident.Id) (meta.Generic, bool)
	// setValue can error, ex. classes.
	setValue(slot ident.Id, value meta.Generic) error
}

func makeProperty(mdl *Metal, p *M.PropertyModel, v valueStore) (ret meta.Property) {
	mem := memProp{mdl, p, v}

	switch t := p.Type; t {
	case M.NumProperty:
		if !p.IsMany {
			ret = &numProp{mem}
		} else {
			ret = &numListProp{mem}
		}

	case M.TextProperty:
		if !p.IsMany {
			ret = &textProp{mem}
		} else {
			ret = &textListProp{mem}
		}
	case M.EnumProperty:
		ret = &enumProp{mem}

	case M.PointerProperty:
		if !p.IsMany {
			// since we dont support "many to many" the far side here must be either "one" or "many".
			// "many" values are "views" -- they inspect the values of this side.
			// "one" values are bi-connections 00 the far side needs to change when this side changes.
			if rel, ok := mdl.Relations[p.Relation]; ok && rel.Style == M.OneToOne {
				ret = &oneToOneProp{mem, rel.GetOther(p.Id)}
			} else {
				ret = &pointerProp{mem}
			}
		} else {
			if rel, ok := mdl.Relations[p.Relation]; ok {
				return &manyToOneProp{mem, rel}
			} else {
				ret = &pointerListProp{mem}
			}
		}
	default:
		panic(errutil.New("property", p.Id, "has unknown type", t))
	}
	return
}

// MARS: it may be necessary to implement individual Gets, especially for instance values, so that we can convert from de-serialized storage to proper return value. for instance, json serializes ident.Id as string. ( keeping in mind that evals break default json serialization anyway. )

type memProp struct {
	mdl   *Metal
	prop  *M.PropertyModel
	value valueStore
}

// String implements Stringer ( for logging )
func (p *memProp) String() string {
	return string(p.prop.Id)
}

// GetId implements meta.Property.
func (p *memProp) GetId() ident.Id {
	return p.prop.Id
}

// GetName implements meta.Property.
func (p *memProp) GetName() string {
	return p.prop.Name
}

// GetRelative implements meta.Property.
// FIX: this exists for client backwards compatiblity.
// the reality is, a relation may have multiple views that need updating. client pulling the views it needs might be best, plus or minus the fact that relation status querried at the end of a frame might not match the relation status recorded by an event.
func (p *memProp) GetRelative() (ret meta.Relative, okay bool) {
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

// SetGeneric panics if not overriden.
func (p *memProp) SetGeneric(value meta.Generic) error {
	panic(errutil.New("set generic not implemented for", p.prop.Id))
}

// SetRelation panics if not overriden.
func (p *memProp) SetRelation(ident.Id) error {
	panic(errutil.New("set relation not supported for property", p.prop.Id))
}

func (p *memProp) getValue() (meta.Generic, bool) {
	return p.value.getValue(p.prop.Id)
}

func (p *memProp) setValue(value meta.Generic) error {
	return p.value.setValue(p.prop.Id, value)
}

type numProp struct{ memProp }

func (p *numProp) GetType() meta.PropertyType {
	return meta.NumProperty
}

func (p *numProp) GetGeneric() (ret meta.Generic) {
	if v, ok := p.value.getValue(p.prop.Id); ok {
		ret = v
	} else {
		var zero rt.Number
		ret = zero
	}
	return
}

// SetGeneric num expects rt.Number
func (p *numProp) SetGeneric(value meta.Generic) (err error) {
	if _, ok := value.(rt.Number); !ok {
		err = errutil.New("set property", p, "expected num, got", sbuf.Type{value})
	} else {
		err = p.setValue(value)
	}
	return
}

type textProp struct{ memProp }

func (p *textProp) GetType() meta.PropertyType {
	return meta.TextProperty
}

// caller can verify property is of expected type if desired;
// we dont need to do that as well.
func (p *textProp) GetGeneric() (ret meta.Generic) {
	if v, ok := p.value.getValue(p.prop.Id); ok {
		ret = v
	} else {
		var zero rt.Text
		ret = zero
	}
	return
}

// SetGeneric text expects rt.Text
func (p *textProp) SetGeneric(value meta.Generic) (err error) {
	if _, ok := value.(rt.Text); !ok {
		err = errutil.New("set property", p, "expected text, got", sbuf.Type{value})
	} else {
		err = p.setValue(value)
	}
	return
}

type enumProp struct{ memProp }

func (p *enumProp) GetType() meta.PropertyType {
	return meta.StateProperty
}

// GetGeneric for enumerated properties returns rt.State.
func (p *enumProp) GetGeneric() (ret meta.Generic) {
	// MARS: this manufactures the state (eval) from the id.
	// would it be better to synthesize this, or even just store it as state in the first place?
	if v, ok := p.value.getValue(p.prop.Id); ok {
		if id, ok := v.(ident.Id); !ok {
			panic(errutil.New("stored enum has invalid value", sbuf.Type{v}))
		} else {
			ret = rt.State(id)
		}
	} else {
		enum := p.mdl.Enumerations[p.prop.Id]
		ret = rt.State(enum.Best())
	}
	return
}

// SetGeneric enum expects rt.State
func (p *enumProp) SetGeneric(value meta.Generic) (err error) {
	if state, ok := value.(rt.State); !ok {
		err = errutil.New("set property", p, "expected state, got", sbuf.Type{value})
	} else {
		// MARS: FIX: constraints!!!
		strip := ident.Id(state)
		err = p.setValue(strip)
	}
	return
}

type pointerProp struct{ memProp }

func (p *pointerProp) GetType() meta.PropertyType {
	return meta.ObjectProperty
}

// caller can verify property is of expected type if desired;
// we dont need to do that as well.
func (p *pointerProp) GetGeneric() (ret meta.Generic) {
	if v, ok := p.value.getValue(p.prop.Id); ok {
		// MARS FIX - pointers are mostly stored as refs, but not relations.
		if id, ok := v.(ident.Id); ok {
			ret = rt.Reference(id)
		} else {
			ret = v
		}
	} else {
		var zero rt.Object // object implements rt.ObjEval
		ret = zero
	}
	return
}

// SetGeneric pointer expects rt.Object, noting rt.Object can be empty.
func (p *pointerProp) SetGeneric(value meta.Generic) (err error) {
	if obj, ok := value.(rt.Object); !ok {
		err = errutil.New("set property", p, "expected object, got", sbuf.Type{value})
	} else if ok := (obj.Instance == nil) || p.mdl.AreCompatible(obj.GetParentClass(), p.prop.Relates); !ok {
		err = errutil.New("set property", p, "object", obj, "not compatible with", p.prop.Relates)
	} else {
		// Objects cant be stored directly, but we can store references to them.
		ref := rt.Reference(obj.GetId())
		err = p.setValue(ref)
	}
	return
}

type numListProp struct{ memProp }

func (p *numListProp) GetType() meta.PropertyType {
	return meta.NumProperty | meta.ArrayProperty
}

func (p *numListProp) GetGeneric() (ret meta.Generic) {
	if v, ok := p.value.getValue(p.prop.Id); ok {
		ret = v
	} else {
		var zero rt.Numbers
		ret = zero
	}
	return
}

// SetGeneric number list expects []rt.Number ( or rt.Number )
func (p *numListProp) SetGeneric(value meta.Generic) (err error) {
	if _, ok := value.([]rt.Number); !ok {
		err = errutil.New("set property", p, "expected numbers, got", sbuf.Type{value})
	} else {
		err = p.setValue(value)
	}
	return
}

type textListProp struct{ memProp }

func (p *textListProp) GetType() meta.PropertyType {
	return meta.TextProperty | meta.ArrayProperty
}

// caller can verify property is of expected type if desired;
// we dont need to do that as well.
func (p *textListProp) GetGeneric() (ret meta.Generic) {
	if v, ok := p.value.getValue(p.prop.Id); ok {
		ret = v
	} else {
		var zero rt.Texts
		ret = zero
	}
	return
}

// SetGeneric text list expects []rt.Text ( or rt.Texts )
func (p *textListProp) SetGeneric(value meta.Generic) (err error) {
	if _, ok := value.([]rt.Text); !ok {
		err = errutil.New("set property", p, "expected texts, got", sbuf.Type{value})
	} else {
		err = p.setValue(value)
	}
	return
}

type pointerListProp struct{ memProp }

func (p *pointerListProp) GetType() meta.PropertyType {
	return meta.ObjectProperty | meta.ArrayProperty
}

func (p *pointerListProp) GetGeneric() (ret meta.Generic) {
	// caller can verify property is of expected type if desired;
	// we dont need to do that as well.
	if v, ok := p.value.getValue(p.prop.Id); ok {
		ret = v
	} else {
		var zero rt.References
		ret = zero
	}
	return
}

// SetGeneric pointer expects []rt.Object ( or rt.Objects )
// MARS - should this be rt.References instead?
func (p *pointerListProp) SetGeneric(value meta.Generic) (err error) {
	if objs, ok := value.([]rt.Object); !ok {
		err = errutil.New("set property", p, "expected objects, got", sbuf.Type{value})
	} else {
		refs := make([]rt.Reference, len(objs))
		for i, obj := range objs {
			if ok := obj.Empty() || p.mdl.AreCompatible(obj.GetParentClass(), p.prop.Relates); !ok {
				err = errutil.New("set property", p, "object", i, obj, "not compatible with", p.prop.Relates)
				break
			}
			refs[i] = rt.Reference(obj.GetId())
		}
		if err == nil {
			err = p.setValue(refs)
		}
	}
	return
}
