package model

import (
	"fmt"
)

//
// Script instances.
// operates kind of like a prototype
// the values have references to the class table
// its properties fallback to the table when needed.
//
type InstanceInfo struct {
	id    StringId
	class *ClassInfo
	name  string
	long  string // FIX: kill this, replace with article categorization
	refs  *References
	enum  map[StringId]int
	text  map[StringId]string
	num   map[StringId]float32
}

//
//
//
func NewInstanceInfo(
	id StringId,
	class *ClassInfo,
	name string,
	long string,
	refs *References,
) *InstanceInfo {
	inst := &InstanceInfo{id, class, name, long, refs,
		make(map[StringId]int),
		make(map[StringId]string),
		make(map[StringId]float32),
	}
	return inst
}

//
//
//
func (this *InstanceInfo) Id() StringId {
	return this.id
}

//
//
//
func (this *InstanceInfo) String() string {
	// FIX: this looks silly when singular starts with a vowel.
	return fmt.Sprintf("%s ( %s: %s )", this.long, this.id, this.class.singular)
}

//
//
//
func (this *InstanceInfo) Name() string {
	return this.name
}

//
//
//
func (this *InstanceInfo) FullName() string {
	name := this.long
	if name == "" {
		name = this.name
	}
	return name
}

//
//
//
func (this *InstanceInfo) Class() *ClassInfo {
	return this.class
}

//
//
//
func (this *InstanceInfo) CompatibleWith(cls *ClassInfo) bool {
	return this.class == cls || this.class.HasParent(cls)
}

//
// return a interface representing the contents of the passed property name
//
func (this *InstanceInfo) ValueByName(name string) (ret IValue, okay bool) {
	if prop, ok := this.class.FindProperty(name); ok {
		switch prop := prop.(type) {
		case *RelativeProperty:
			ret = &RelativeValue{this, prop}
			okay = true
		case *TextProperty:
			ret = &TextValue{this, prop}
			okay = true
		case *EnumProperty:
			ret = &EnumValue{this, prop, nil}
			okay = true
		case *NumProperty:
			ret = &NumValue{this, prop}
			okay = true
		default:
			panic(fmt.Sprintf("unhandled property %s type %T", name, prop))
		}
	}
	return ret, okay
}

//
//
//
func (this *InstanceInfo) EnumValue(name string) (ret *EnumValue, okay bool) {
	if prop, ok := this.class.FindProperty(name); ok {
		if prop, ok := prop.(*EnumProperty); ok {
			ret = &EnumValue{this, prop, nil}
			okay = ok
		}
	}
	return ret, okay
}

//
//
//
func (this *InstanceInfo) NumValue(name string) (ret *NumValue, okay bool) {
	if prop, ok := this.class.FindProperty(name); ok {
		if prop, ok := prop.(*NumProperty); ok {
			ret = &NumValue{this, prop}
			okay = ok
		}
	}
	return ret, okay
}

//
//
//
func (this *InstanceInfo) TextValue(name string) (ret *TextValue, okay bool) {
	if prop, ok := this.class.FindProperty(name); ok {
		if prop, ok := prop.(*TextProperty); ok {
			ret = &TextValue{this, prop}
			okay = ok
		}
	}
	return ret, okay
}
