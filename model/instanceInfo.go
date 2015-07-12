package model

import (
	"fmt"
	"github.com/ionous/sashimi/util/ident"
)

//
// Script Instances operate kind of like a prototype:
// its property values fall back to its associated class when needed.
//
type InstanceInfo struct {
	id     ident.Id
	class  *ClassInfo
	name   string // FIX: kill name, replace wth name text value.
	long   string // FIX: kill long name, replace with article categorization. somehow.
	values map[ident.Id]Variant
}

//
// Enums are stored as int;
// Numbers as float32;
// Pointers as ident.Id;
// Text as string.
//
type Variant interface{}

//
// NewInstanceInfo creates a new instance record.
//
func NewInstanceInfo(
	id ident.Id,
	class *ClassInfo,
	name string,
	long string,
	values map[ident.Id]Variant,
) *InstanceInfo {
	inst := &InstanceInfo{id, class, name, long, values}
	return inst
}

//
// Id to uniquely identify the instance.
// Usually, the id is derived automatically from the instance's name.
//
func (inst *InstanceInfo) Id() ident.Id {
	return inst.id
}

//
// String representation of the instance and its class.
//
func (inst *InstanceInfo) String() string {
	// FIX: inst looks silly when singular starts with a vowel.
	return fmt.Sprintf("%s ( %s: %s )", inst.long, inst.id, inst.class.singular)
}

//
// Name given by the user for this instance ( sans articles of speech. )
//
func (inst *InstanceInfo) Name() string {
	return inst.name
}

//
// Class of this instance.
//
func (inst *InstanceInfo) Class() *ClassInfo {
	return inst.class
}

//
// FindValue returns the (default) value of the named property.
// By design, there is no equivalent SetValue.
//
func (inst *InstanceInfo) FindValue(name string) (ret interface{}, okay bool) {
	if prop, ok := inst.class.FindProperty(name); ok {
		ret, okay = inst._value(prop), true
	}
	return ret, okay
}

func (inst *InstanceInfo) Values() map[ident.Id]Variant {
	return inst.values
}

//
// Value returns the (default) value of the identified property.
// By design, there is no equivalent SetValue.
//
func (inst *InstanceInfo) Value(id ident.Id) (ret interface{}, okay bool) {
	if prop, ok := inst.class.PropertyById(id); ok {
		ret, okay = inst._value(prop), true
	}
	return ret, okay
}

func (inst *InstanceInfo) _value(prop IProperty) (ret interface{}) {
	if v, ok := inst.values[prop.Id()]; ok {
		ret = v
	} else {
		ret = prop.Zero(inst.class.Constraints())
	}
	return ret
}
