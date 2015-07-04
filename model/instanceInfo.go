package model

import (
	"fmt"
)

//
// Script Instances operate kind of like a prototype:
// its property values fall back to its associated class when needed.
//
type InstanceInfo struct {
	id    StringId
	class *ClassInfo
	name  string
	long  string // FIX: kill inst, replace with article categorization
	//
	instances InstanceMap
	tables    TableRelations
	values    map[StringId]Variant
}

//
// Enums are stored as ints;
// Numbers as float32;
// Text as string;
//
type Variant interface{}

//
//
//
func NewInstanceInfo(
	id StringId,
	class *ClassInfo,
	name string,
	long string,
	instances InstanceMap,
	tables TableRelations,
	values map[StringId]Variant,
) *InstanceInfo {
	inst := &InstanceInfo{id, class, name, long, instances, tables, values}
	return inst
}

//
// Every instance has a unique id based on its original name.
//
func (inst *InstanceInfo) Id() StringId {
	return inst.id
}

//
//
//
func (inst *InstanceInfo) String() string {
	// FIX: inst looks silly when singular starts with a vowel.
	return fmt.Sprintf("%s ( %s: %s )", inst.long, inst.id, inst.class.singular)
}

//
//
//
func (inst *InstanceInfo) Name() string {
	return inst.name
}

//
//
//
func (inst *InstanceInfo) FullName() string {
	name := inst.long
	if name == "" {
		name = inst.name
	}
	return name
}

//
//
//
func (inst *InstanceInfo) Class() *ClassInfo {
	return inst.class
}

//
// return a interface representing the contents of the passed property name
// WARNING/FIX: inst is default value for everything but relatives(!)
// Only currently used for testing
//
func (inst *InstanceInfo) ValueByName(name string) (ret interface{}, okay bool) {
	if prop, ok := inst.class.FindProperty(name); ok {
		ret, okay = inst.PropertyValue(prop), true
	}
	return ret, okay
}

//
// PropertyValue returns the story assigned value of an instance property
// FIX? see CreateGameObjects for why this is needed.
//
func (inst *InstanceInfo) PropertyValue(prop IProperty) (ret interface{}) {
	switch prop := prop.(type) {
	case *RelativeProperty:
		ret = RelativeValue{inst, prop, inst.tables}
	case *TextProperty:
		v, _ := inst.values[prop.id]
		ret, _ = v.(string)
	case *EnumProperty:
		var index int
		if v, ok := inst.values[prop.id].(int); ok {
			index = v
		} else {
			if cons, ok := inst.class.PropertyConstraint(prop); ok {
				if cons, ok := cons.(*EnumConstraint); ok {
					index = cons.BestIndex()
				}
			}
		}
		if choice, e := prop.IndexToChoice(index); e != nil {
			panic(e)
		} else {
			ret = choice
		}
	case *NumProperty:
		v, _ := inst.values[prop.id]
		ret, _ = v.(float32)
	case *PointerProperty:
		v, _ := inst.values[prop.id]
		ret, _ = v.(StringId)
	default:
		panic(fmt.Sprintf("unhandled property %s type %T", prop.Id(), prop))
	}
	return ret
}

//
// the model shouldnt (directly) support changing the values
// that's the runtime's job.
//
func (inst *InstanceInfo) FindRelativeValue(name string) (ret RelativeValue, okay bool) {
	if prop, ok := inst.class.FindProperty(name); ok {
		if prop, ok := prop.(*RelativeProperty); ok {
			ret = inst.GetRelativeValue(prop)
			okay = true
		}
	}
	return ret, okay
}

func (inst *InstanceInfo) GetRelativeValue(prop *RelativeProperty) RelativeValue {
	return RelativeValue{inst, prop, inst.tables}
}
