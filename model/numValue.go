package model

import "fmt"

//
// Number Values
//
type NumValue struct {
	inst *InstanceInfo
	prop *NumProperty
}

//
func (this *NumValue) Property() IProperty {
	return this.prop
}

//
func (this *NumValue) Any() (interface{}, bool) {
	return this.Num()
}

//
func (this *NumValue) String() string {
	value, _ := this.Num()
	return fmt.Sprint(value)
}

//
func (this *NumValue) Num() (float32, bool) {
	num, ok := this.inst.values[this.prop.id].(float32)
	return num, ok
}
