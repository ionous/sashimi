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
func (this *NumValue) Any() interface{} {
	return this.Num()
}

//
func (this *NumValue) String() string {
	return fmt.Sprint(this.Num())
}

//
func (this *NumValue) Num() float32 {
	num, _ := this.inst.values[this.prop.id].(float32)
	return num
}
