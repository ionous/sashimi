package model

//
// Text variable used by game objects.
//
type TextValue struct {
	inst *InstanceInfo
	prop *TextProperty
}

//
func (this *TextValue) Property() IProperty {
	return this.prop
}

//
func (this *TextValue) Any() interface{} {
	return this.Text()
}

//
func (this *TextValue) String() string {
	return this.Text()
}

//
func (this *TextValue) Text() string {
	text, _ := this.inst.values[this.prop.id].(string)
	return text
}
