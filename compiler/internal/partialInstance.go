package internal

import (
	M "github.com/ionous/sashimi/compiler/xmodel"
	"github.com/ionous/sashimi/util/ident"
)

// PartialInstance binds an instance to one particular class,
// and begins to fill out the values for the future object.
type PartialInstance struct {
	id               ident.Id // unique id for the instance
	name             string
	longName         string
	class            *M.ClassInfo     // the class this object will use to construct itself
	propertyBuilders PropertyBuilders // access to the class data
	values           PendingValues    // accumulates the object's initial values
	refs             PartialMap       // verification for the existance of other instances
	tables           TableRelations   // access to relation data
}

// setKeyValue, helper to set instance property choices.
func (inst *PartialInstance) setKeyValue(id ident.Id, value interface{}) (err error) {
	if builder, ok := inst.propertyBuilders.propertyById(id); !ok {
		err = PropertyNotFound(inst.class.Id, id.String())
	} else {
		err = inst.setProperty(builder, value)
	}
	return err
}

// setNameValue, helper to set instance property values.
func (inst *PartialInstance) setNameValue(name string, value interface{}) (err error) {
	if builder, ok := inst.propertyBuilders.findProperty(name); !ok {
		err = PropertyNotFound(inst.class.Id, name)
	} else {
		err = inst.setProperty(builder, value)
	}
	return err
}

func (inst *PartialInstance) setProperty(builder IBuildProperty, value interface{}) (err error) {
	return builder.SetProperty(PropertyContext{
		inst.id,
		inst.tables,
		inst.class,
		inst.values,
		inst.refs,
		value,
	})
}
