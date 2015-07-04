package model

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
func (this *PointerValue) Any() interface{} {
	return this.Pointer()
}

//
func (this *PointerValue) String() string {
	return this.Pointer().String()
}

//
func (this *PointerValue) Pointer() StringId {
	text, _ := this.inst.values[this.prop.id].(StringId)
	return text
}
