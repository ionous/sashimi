package model

import "fmt"

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
	text, ok := this.inst.text[this.prop.id]
	return text, ok
}

//
func (this *TextValue) SetText(text string) {
	this.inst.text[this.prop.id] = text
}

//
func (this *TextValue) SetAny(value interface{}) (err error) {
	if text, ok := value.(string); !ok {
		err = fmt.Errorf(" %#v set any failed with %#v", this, value)
	} else {
		this.SetText(text)
	}
	return err
}
