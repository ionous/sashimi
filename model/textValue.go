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
func (this *TextValue) Any() (interface{}, bool) {
	return this.Text()
}

//
func (this *TextValue) String() string {
	text, _ := this.Text()
	return text
}

//
func (this *TextValue) Text() (string, bool) {
	text, ok := this.inst.values[this.prop.id].(string)
	return text, ok
}
