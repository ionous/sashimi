package metal

import (
	"github.com/ionous/mars/rt"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/errutil"
	"github.com/ionous/sashimi/util/ident"
	"github.com/ionous/sashimi/util/sbuf"
)

// in sashimi v1 was newRelatedValue()->objectWriteValue
// FIX? i think the synchronization between sides of the property might be better done as a property watcher.
// it could live on a whole other layer, and the relation could drive the watch,
// and the Relation value in the PropertyModel wouldnt be needed.
type oneToOneProp struct {
	memProp
	targetProp ident.Id
}

func (p *oneToOneProp) GetType() meta.PropertyType {
	return meta.ObjectProperty | meta.ArrayProperty
}

// GetGeneric for one-to-one relations returns an object eval.
func (p *oneToOneProp) GetGeneric() (ret meta.Generic) {
	if id, e := p.getOther(); e != nil {
		panic(errutil.New("stored one-to-one has invalid value", e))
	} else {
		// FIX: why cant we store references directly
		ret = rt.Reference(id)
	}
	return
}

func (p *oneToOneProp) SetGeneric(value meta.Generic) (err error) {
	if obj, ok := value.(rt.Object); !ok {
		err = errutil.New("one-to-one", p, "expected object, got", sbuf.Type{value})
	} else if wasId, e := p.getOther(); e != nil {
		err = e
	} else if newId := obj.GetId(); newId != wasId {
		if !obj.Exists() {
			// clearing value? this is easy
			p.setValue(ident.Empty())
		} else {
			// changing to something new:
			if e := p.areCompatible(obj.GetParentClass()); e != nil {
				err = e // ^^ this cant store that
			} else if far, e := p.getFarProp(obj); e != nil {
				err = e // ^^ that should have been a one-to-one-property
			} else if e := far.areCompatible(p.value.getClassId()); e != nil {
				err = e // ^^ that cant store this
			} else if e := p.setValue(value); e != nil {
				err = errutil.New("one-to-one", p, "error seting value", e)
			} else if e := far.setValue(p.value.getStoreId()); e != nil {
				err = e // ^^ an error has occured, but lets try to roll back.
				if e := p.setValue(wasId); e != nil {
					err = errutil.New("one-to-one", p, "couldnt roll back value", err, e)
				}
			}
		}
		// clear old value:
		if !wasId.Empty() {
			old, _ := p.mdl.getInstance(wasId)
			if other, e := p.getFarProp(old); e != nil {
				err = errutil.New("one-to-one", p, "error getting old value", e)
			} else if e := other.setValue(ident.Empty()); e != nil {
				err = errutil.New("one-to-one", p, "error clearing old value", e)
			}

		}
	}
	return
}

// areCompatible returns true if the passed class can be used by this relation.
func (p *oneToOneProp) areCompatible(class ident.Id) (err error) {
	if ok := p.mdl.AreCompatible(class, p.prop.Relates); !ok {
		err = errutil.New("one-to-one", p, "class", class, "not compatible with", p.prop.Relates)
	}
	return
}

// getOther returns the id of the thing we point to
// MARS: relations are stored by id; is that right?
func (p *oneToOneProp) getOther() (ret ident.Id, err error) {
	// if the value doesnt exist, we use zero.
	if v, ok := p.getValue(); ok {
		if id, ok := v.(ident.Id); !ok {
			err = errutil.New("relation has unexpected storage type", sbuf.Type{v})
		} else {
			ret = id
		}
	}
	return
}

// getFarProp returns the reverse one-to-one property from passed instance
func (p *oneToOneProp) getFarProp(other meta.Instance) (ret *oneToOneProp, err error) {
	if other == nil {
		err = errutil.New("one-to-one", p, "missing far object", other)
	} else if prop, ok := other.GetProperty(p.targetProp); !ok {
		err = errutil.New("one-to-one", p, "missing target", p.targetProp)
	} else if other, ok := prop.(*oneToOneProp); !ok {
		err = errutil.New("one-to-one", p, "mismatched target", p.targetProp, sbuf.Type{other})
	} else {
		ret = other
	}
	return
}
