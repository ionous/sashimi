package compiler

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
	S "github.com/ionous/sashimi/source"
	"github.com/ionous/sashimi/util/errutil"
	"log"
)

type PartialInstances struct {
	tables   M.TableRelations
	partials PartialMap
}
type PartialMap map[M.StringId]*PartialInstance

func newPartialInstances(relations M.RelationMap) PartialInstances {
	tables := M.NewTableRelations(relations)
	partials := make(PartialMap)
	return PartialInstances{tables, partials}
}

//
// Add a new set of references for the passed id'd reference.
//
func (part *PartialInstances) newInstance(pending *PendingInstance, class *M.ClassInfo, props PropertyBuilders) {
	id, name, long := pending.id, pending.name, pending.longName
	values := make(PendingValues)
	partial := &PartialInstance{id, name, long, class, props, values, part.partials, part.tables}
	part.partials[id] = partial
}

//
// makeData sets all pending data to the known instances.
// returns thos instance and the tables
//
func (part *PartialInstances) makeData(choices []S.ChoiceStatement, kvs []S.KeyValueStatement,
) (instances M.InstanceMap, tables M.TableRelations, err error) {
	if e := part._addChoices(choices); e != nil {
		err = e
	} else if e := part._addKeyValues(kvs); e != nil {
		err = e
	}
	tables = part.tables
	instances = make(M.InstanceMap)
	for id, p := range part.partials {
		instance := M.NewInstanceInfo(p.id, p.class, p.name, p.longName, instances, tables, p.values)
		instances[id] = instance
	}

	return instances, tables, err
}

//
// via makeData(): Add key value data to the targeted instances
//
func (part *PartialInstances) _addChoices(choices []S.ChoiceStatement) (err error) {
	log.Println("adding instance choices")
	for _, choice := range choices {
		fields := choice.Fields()
		id := M.MakeStringId(fields.Owner)
		if inst, ok := part.partials[id]; !ok {
			err = errutil.Append(err, M.InstanceNotFound(fields.Owner))
		} else if prop, index, ok := inst.class.PropertyByChoice(fields.Choice); !ok {
			e := fmt.Errorf("no such choice: '%v'", choice)
			err = errutil.Append(err, e)
		} else if e := inst.setKeyValue(prop.Name(), index); e != nil {
			err = errutil.Append(err, e)
		}

	}
	return err
}

//
// via makeData(): Add key value data to the targeted instances
//
func (part *PartialInstances) _addKeyValues(kvs []S.KeyValueStatement) (err error) {
	log.Println("adding instance key values")
	for _, kv := range kvs {
		fields := kv.Fields()
		id := M.MakeStringId(fields.Owner)
		if inst, ok := part.partials[id]; !ok {
			err = errutil.Append(err, M.InstanceNotFound(fields.Owner))
		} else if e := inst.setKeyValue(fields.Key, fields.Value); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return err
}
