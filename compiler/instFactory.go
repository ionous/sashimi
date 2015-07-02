package compiler

import (
	M "github.com/ionous/sashimi/model"
	S "github.com/ionous/sashimi/source"
	"github.com/ionous/sashimi/util/errutil"
)

type InstanceFactory struct {
	names   NameSource
	pending PendingInstances
}
type PendingInstances map[M.StringId]*PendingInstance

func newInstanceFactory(names NameSource) *InstanceFactory {
	return &InstanceFactory{names, make(PendingInstances)}
}

func (this *InstanceFactory) findPendingInstance(name string,
) (*PendingInstance, bool,
) {
	id, _ := this.names.addName(nil, name, "instance", "")
	ret, okay := this.pending[id]
	return ret, okay
}

//
// Register the passed `name` as an instance of `class`
// NOTE: there can be multiple assertions refering to the same instance.
//
func (this *InstanceFactory) addInstanceRef(inst string, cls M.StringId, options S.Options, code S.Code,
) (ret *PendingInstance, err error,
) {
	id, err := this.names.addName(nil, inst, "instance", "")
	if i, ok := this.pending[id]; ok {
		ret = i
	} else {
		ret = &PendingInstance{id: id, name: inst}
		this.pending[id] = ret
	}
	ret.classes.addClassReference(cls, code)
	//1
	longName := options["long name"]
	if longName != "" {
		ret.longName = longName
	}
	return ret, err
}

//
func (this *InstanceFactory) makeInstances(log *ErrorLog, classes M.ClassMap, relations M.RelationMap) (
	ret PartialInstances,
	err error,
) {
	inner := make(M.InstanceMap)
	tables := M.NewTableRelations(relations)
	refs := M.NewReferences(classes, inner, tables)
	instances := PartialInstances{log, inner, tables}
	// resolve all of the classes and create instances for them
	for id, pending := range this.pending {
		if class, e := pending.classes.resolveClass(classes); e != nil {
			err = errutil.Append(err, e)
		} else {
			// NOTE: this implicitly adds to the inner instances list
			refs.NewInstance(id, class, pending.name, pending.longName)
		}
	}
	return instances, err
}
