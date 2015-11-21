package metal

import (
	"fmt"
	//	M "github.com/ionous/sashimi/compiler/model"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/errutil"
	"github.com/ionous/sashimi/util/ident"
)

var _ = fmt.Println

type objectList struct {
	panicValue
	targetProp ident.Id
	objs       []ident.Id
}

// the many side of a many-to-one, or one-to-many relation;
// returns a list.
func newManyValues(p *propBase) (ret meta.Values) {
	var objs []ident.Id
	rel, ok := p.mdl.Relations[p.prop.Relation]
	if !ok {
		panic(fmt.Sprintf("missing relation '%s'", p.prop.Relation))
	}
	// check instance because newManyValues can be called by class values ( getZero )
	var targetProp ident.Id
	if _, ok := p.mdl.Instances[p.src]; ok {
		// FIX: would rather make this a datastore query;
		// ( would require changing from ObjectValue interface to a full shadow model. )
		targetProp = rel.GetOther(p.prop.Id)
		// use the meta interface in order to get latest data
		for i := 0; i < p.mdl.NumInstance(); i++ {
			target := p.mdl.InstanceNum(i)
			if t, ok := target.GetProperty(targetProp); ok {
				if v := t.GetValue(); v.GetObject() == p.src {
					objs = append(objs, target.GetId())
				}
			}
		}
	}
	return &objectList{panicValue{p}, targetProp, objs}
}

func (p objectList) NumValue() int {
	return len(p.objs)
}

func (p objectList) ValueNum(i int) meta.Value {
	return objectReadValue{p.panicValue, p.objs[i]}
}

func (p *objectList) ClearValues() (err error) {
	for _, id := range p.objs {
		if v, e := p.mdl.getFarPointer(id, p.targetProp); e != nil {
			err = errutil.Append(err, e)
		} else {
			// possible, if unlikely, that its changed.
			if v.GetObject() == p.src {
				v.SetObject(ident.Empty())
			}
		}
	}
	p.objs = nil
	return
}

// AppendObject writes *this* object into the destination
// ( and updates our list of objects )
func (p *objectList) AppendObject(id ident.Id) (err error) {
	// have to use the meta interface in order to trigger shadow properties
	if v, e := p.mdl.getFarPointer(id, p.targetProp); e != nil {
		err = e
	} else {
		found := false
		for _, already := range p.objs {
			if already == id {
				found = true
				break
			}
		}
		if !found {
			if e := v.SetObject(p.src); e != nil {
				err = errutil.Append(err, e)
			} else {
				p.objs = append(p.objs, id)
			}
		}
	}
	return
}
