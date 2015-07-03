package model

import "fmt"

//
// PointerValue points from one instance to another.
//
type PointerValue struct {
	inst *InstanceInfo
	prop *PointerProperty
}

//
func (this *PointerValue) Property() IProperty {
	return this.prop
}

//
func (this *PointerValue) Any() (interface{}, bool) {
	return this.Pointer()
}

//
func (this *PointerValue) String() string {
	id, _ := this.Pointer()
	return id.String()
}

//
func (this *PointerValue) Pointer() (StringId, bool) {
	text, ok := this.inst.values[this.prop.id].(StringId)
	return text, ok
}

//
func (this *PointerValue) SetPointer(id StringId) {
	this.inst.values[this.prop.id] = id
}

//
func (this *PointerValue) SetAny(value interface{}) (err error) {
	if id, ok := value.(StringId); !ok {
		err = fmt.Errorf(" %#v set any failed with %#v", this, value)
	} else {
		this.SetPointer(id)
	}
	return err
}
