package metal

import (
	"github.com/ionous/mars/rt"
	M "github.com/ionous/sashimi/compiler/model"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/errutil"
	"github.com/ionous/sashimi/util/ident"
)

type manyToOneProp struct {
	memProp
	rel *M.RelationModel
}

func (p *manyToOneProp) GetType() meta.PropertyType {
	return meta.ObjectProperty | meta.ArrayProperty
}

// GetGeneric for one-to-one relations returns an object list eval
func (p *manyToOneProp) GetGeneric() meta.Generic {
	return p
}

func (p *manyToOneProp) GetObjStream(run rt.Runtime) (ret rt.ObjectStream, err error) {
	if rel, ok := p.mdl.Relations[p.prop.Relation]; !ok {
		err = errutil.New("missing relation", p.prop.Relation)
	} else {
		targetProp := rel.GetOther(p.prop.Id)
		it := &manyIt{p, targetProp, p.mdl.NumInstance(), nil}
		ret, it.next = it, it.advance()
	}
	return
}

// SetGeneric for many-to-on relations is invalid.
func (p *manyToOneProp) SetGeneric(value meta.Generic) error {
	return errutil.New("you cant write to many-to-one relations")
}

type manyIt struct {
	p          *manyToOneProp
	targetProp ident.Id
	idx        int
	next       *rt.Object
}

func (it *manyIt) HasNext() bool {
	return it.next != nil
}

func (it *manyIt) GetNext() (ret rt.Object, err error) {
	if it.next == nil {
		err = errutil.New("out of range")
	} else {
		ret, it.next = *it.next, it.advance()
	}
	return
}

func (it *manyIt) advance() (ret *rt.Object) {
	// classes have properties too so check for an empty storage id
	// MARS: fix, what about default relations? maybe no such thing.
	if myId := it.p.value.getStoreId(); !myId.Empty() {
		// advance by dec so we dont have to keep asking for instance count.
		for it.idx > 0 {
			it.idx--
			target := it.p.mdl.InstanceNum(it.idx)
			// not every instance will have the property
			if t, ok := target.GetProperty(it.targetProp); ok {
				// now that we have the property, get its value;
				// i think we are storing idents for relations
				v := t.GetGeneric()
				if mightBeMe, ok := v.(ident.Id); !ok {
					// FIX!!!! pointer prop now returns reference re: mars fix.
					//panic("stored one-to-many has invalid value")
					// haave to change to use underlying instance interface
					panic(v)
				} else if mightBeMe.Equals(myId) {
					// while we've got it: wrap up a direct pointer
					ret = &rt.Object{target}
					break
				}
			}
		}
	}
	return
}
