package compiler

import (
	M "github.com/ionous/sashimi/model"
	S "github.com/ionous/sashimi/source"
)

type InstanceFactory struct {
	names   NameSource
	pending PendingInstances
}
type PendingInstances map[M.StringId]*PendingInstance

func newInstanceFactory(names NameSource) *InstanceFactory {
	return &InstanceFactory{names, make(PendingInstances)}
}

//
// register the passed `name` as an instance of `class`
// NOTE: there can be multiple assertions refering to the same instance.
func (this *InstanceFactory) addInstanceRef(class ClassReference, name string, options S.Options,
) (inst *PendingInstance, err error,
) {
	id, err := this.names.addName(nil, name, "instance")
	if i, ok := this.pending[id]; ok {
		inst = i
	} else {
		inst = &PendingInstance{name: name}
		this.pending[id] = inst
	}
	inst.classes.addClassReference(class)
	//
	longName := options["long name"]
	if longName != "" {
		inst.longName = longName
	}
	return inst, err
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
			err = AppendError(err, e)
		} else {
			// NOTE: this implicitly adds to the inner instances list
			refs.NewInstance(id, class, pending.name, pending.longName)
		}
	}
	return instances, err
}
