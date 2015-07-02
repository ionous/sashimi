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

//
func (this *NumValue) SetNum(num float32) {
	this.inst.values[this.prop.id] = num
}

//
func (this *NumValue) SetAny(value interface{}) (err error) {
	switch num := value.(type) {
	case nil:
		this.SetNum(0)
	case int:
		this.SetNum(float32(num))
	case float32:
		this.SetNum(num)
	case float64: // note: go's own default number type is float64
		this.SetNum(float32(num))
	default:
		err = fmt.Errorf(" %#v set any failed with %#v", this, value)
	}
	return err
}
