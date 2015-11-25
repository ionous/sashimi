package metal

import (
	"fmt"
	"github.com/ionous/sashimi/util/ident"
)

// PanicValue implements the Value(s) interface:
// panicing on every get() and set(), and then
// specific property types override the specific methods they need:
// .text for text, num for num, etc.
type panicValue struct {
	*propBase
}

func (p panicValue) GetNum() float32 {
	panic(fmt.Errorf("get num not supported for property %s", p.prop))
}
func (p panicValue) SetNum(float32) error {
	panic(fmt.Errorf("set num not supported for property %s", p.prop))
}
func (p panicValue) AppendNum(float32) error {
	panic(fmt.Errorf("append num not supported for property %s", p.prop))
}
func (p panicValue) GetText() string {
	panic(fmt.Errorf("get text not supported for property %s", p.prop))
}
func (p panicValue) SetText(string) error {
	panic(fmt.Errorf("set text not supported for property %s", p.prop))
}
func (p panicValue) AppendText(string) error {
	panic(fmt.Errorf("append text not supported for property %s", p.prop))
}
func (p panicValue) GetObject() ident.Id {
	panic(fmt.Errorf("get object not supported for property %s", p.prop))
}
func (p panicValue) SetObject(ident.Id) error {
	panic(fmt.Errorf("set object not supported for property %s", p.prop))
}
func (p panicValue) AppendObject(ident.Id) error {
	panic(fmt.Errorf("append object not supported for property %s", p.prop))
}
func (p panicValue) GetState() ident.Id {
	panic(fmt.Errorf("get state not supported for property %s", p.prop))
}
func (p panicValue) SetState(ident.Id) error {
	panic(fmt.Errorf("set state not supported for property %s", p.prop))
}
