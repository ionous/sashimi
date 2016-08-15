package xmodel

import (
	"fmt"
	"github.com/ionous/sashimi/util/ident"
)

// Script Instances operate kind of like a prototype:
// its property values fall back to its associated class when needed.
type InstanceInfo struct {
	// Id to uniquely identify the instance.
	// Usually, the id is derived automatically from the instance's name.
	Id    ident.Id
	Class *ClassInfo
	//Name given by the user for this instance ( sans articles of speech.
	Name string // FIX: kill name, replace wth name text value.
	// /long   string // FIX: kill long name, replace with article categorization. somehow.
	Values map[ident.Id]Variant
}

// Enums are stored as int;
// Numbers as float64;
// Pointers as ident.Id;
// Text as string.
type Variant interface{}

// String representation of the instance and its class.
func (inst *InstanceInfo) String() string {
	// FIX: inst looks silly when singular starts with a vowel.
	return fmt.Sprintf("%s(%s)", inst.Id, inst.Class.Singular)
}

// FindValue returns the (default) value of the named property.
// By design, there is no equivalent SetValue.
func (inst *InstanceInfo) FindValue(name string) (ret interface{}, okay bool) {
	if prop, ok := inst.Class.FindProperty(name); ok {
		ret, okay = inst._value(prop)
	}
	return
}

// Value returns the (default) value of the identified property.
// By design, there is no equivalent SetValue.
func (inst *InstanceInfo) Value(id ident.Id) (ret interface{}, okay bool) {
	if prop, ok := inst.Class.GetProperty(id); ok {
		ret, okay = inst._value(prop)
	}
	return
}

func (inst *InstanceInfo) _value(prop IProperty) (ret interface{}, okay bool) {
	ret, okay = inst.Values[prop.GetId()]
	return
}
