package metal

import (
	"fmt"
	M "github.com/ionous/sashimi/compiler/model"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
)

// returned by object/relation array b/c we cant mutate individual values
type objectReadValue struct {
	panicValue
	currentVal ident.Id
}

func (p objectReadValue) GetObject() (ret ident.Id) {
	return p.currentVal
}

// the one side of a many-to-one, one-to-one, or one-to-many relation.
func newRelatedValue(p *propBase, rel *M.RelationModel) meta.Value {
	return objectWriteValue{panicValue{p}, rel.GetOther(p.prop.Id)}
}

// returned by newRelatedValue
// FIX? i think the synchronization between sides of the property might be better done as a property watcher.
// it could live on a whole other layer, and the relation could drive the watch,
// and the Relation value in the PropertyModel wouldnt be needed.
type objectWriteValue struct {
	panicValue
	targetProp ident.Id
}

func (p objectWriteValue) GetObject() ident.Id {
	return p.get().(ident.Id)
}

func (p objectWriteValue) SetObject(id ident.Id) (err error) {
	if was := p.GetObject(); was != id {
		// clear other
		if !was.Empty() {
			if v, e := p.mdl.getFarPointer(was, p.targetProp); e != nil {
				err = e
			} else {
				v.SetObject(ident.Empty())
			}
		}
		if err == nil {
			// set this
			if id.Empty() {
				err = p.set(id)
			} else {
				// check that the target object is allowed
				if e := p.mdl.match(id, p.prop.Relates); e != nil {
					err = e
				} else {
					if e := p.set(id); e != nil {
						err = e
					} else {
						// set the reverse,( if the reverse is also a one-to-one )
						if v, e := p.mdl.getFarPointer(id, p.targetProp); e != nil {
							err = e
						} else {
							v.SetObject(p.src)
						}
					}
				}
			}
		}
	}
	return
}

func (mdl *Metal) getFarPointer(target, targetProp ident.Id) (ret meta.Value, err error) {
	if inst, ok := mdl.GetInstance(target); !ok {
		err = fmt.Errorf("couldnt find inst object %s", target)
	} else if prop, ok := inst.GetProperty(targetProp); !ok {
		err = fmt.Errorf("unexpectedly missing target prop %s.%s", target, targetProp)
	} else {
		ret = prop.GetValue()
	}
	return
}
