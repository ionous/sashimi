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
//
func (inst *InstanceInfo) ValueByName(name string) (ret IValue, okay bool) {
	if prop, ok := inst.class.FindProperty(name); ok {
		ret, okay = inst.PropertyValue(prop)
	}
	return ret, okay
}

func (inst *InstanceInfo) PropertyValue(prop IProperty) (ret IValue, okay bool) {
	switch prop := prop.(type) {
	case *RelativeProperty:
		ret = &RelativeValue{inst, prop, inst.tables}
		okay = true
	case *TextProperty:
		ret = &TextValue{inst, prop}
		okay = true
	case *EnumProperty:
		ret = &EnumValue{inst, prop}
		okay = true
	case *NumProperty:
		ret = &NumValue{inst, prop}
		okay = true
	case *PointerProperty:
		ret = &PointerValue{inst, prop}
		okay = true
	default:
		panic(fmt.Sprintf("unhandled property %s type %T", prop.Id(), prop))
	}
	return ret, okay
}

//
// FIX: see ValueByName() and IValue
// the model shouldnt (directly) support changing the values
// that's the runtime's job.
//
func (inst *InstanceInfo) GetRelativeValue(name string) (ret *RelativeValue, okay bool) {
	if prop, ok := inst.class.FindProperty(name); ok {
		if prop, ok := prop.(*RelativeProperty); ok {
			ret = &RelativeValue{inst, prop, inst.tables}
			okay = ok
		}
	}
	return ret, okay
}
