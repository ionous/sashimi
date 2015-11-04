package internal

import (
	M "github.com/ionous/sashimi/model"
	S "github.com/ionous/sashimi/source"
	"github.com/ionous/sashimi/util/errutil"
	"github.com/ionous/sashimi/util/ident"
	"log"
)

type InstanceFactory struct {
	names   NameSource
	log     *log.Logger
	pending PendingInstances
}
type PendingInstances map[ident.Id]*PendingInstance

func NewInstanceFactory(names NameSource, log *log.Logger) *InstanceFactory {
	return &InstanceFactory{names, log, make(PendingInstances)}
}

//
// Register the passed `name` as an instance of `class`
// NOTE: there can be multiple assertions refering to the same instance.
//
func (fact *InstanceFactory) addInstanceRef(cls *PendingClass, name string, longName string, code S.Code,
) (ret *PendingInstance, err error,
) {
	id, err := fact.names.addName(nil, name, "instance", "")
	if i, ok := fact.pending[id]; ok {
		ret = i
	} else {
		ret = &PendingInstance{id: id, name: name}
		fact.pending[id] = ret
	}
	ret.classes.addClassReference(cls, code)
	//1
	if longName != "" {
		ret.longName = longName
	}
	return ret, err
}

//
// MakeInstances resolves all of the classes for the pending instances.
// Returns "partial instances" which allows for setting instance values, and finally baking the model instance.
//
func (fact *InstanceFactory) makeInstances(classes M.ClassMap, relations M.RelationMap) (
	ret PartialInstances,
	err error,
) {
	partials := newPartialInstances(fact.log, relations)
	// resolve all of the classes and create instances for them
	for _, pending := range fact.pending {
		if class, props, e := pending.classes.resolveClass(classes); e != nil {
			err = errutil.Append(err, e)
		} else {
			// NOTE: fact implicitly adds to the inner instances list
			partials.newInstance(pending, class, props)
		}
	}
	return partials, err
}
