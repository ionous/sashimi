package memory

import (
	"fmt"
	"github.com/ionous/sashimi/util/ident"
)

// PanicValue implements the Value interface:
// pancing on every get() and set(), and then
// specific property types override the specific methods they need:
// .text for text, num for num, etc.
type panicValue propBase

func (p panicValue) GetNum() float32 {
	panic(fmt.Errorf("get num not supported for property %v", p.prop.GetId()))
}
func (p panicValue) SetNum(float32) error {
	panic(fmt.Errorf("set num not supported for property %v", p.prop.GetId()))
}
func (p panicValue) GetText() string {
	panic(fmt.Errorf("get text not supported for property %v", p.prop.GetId()))
}
func (p panicValue) SetText(string) error {
	panic(fmt.Errorf("set text not supported for property %v", p.prop.GetId()))
}
func (p panicValue) GetState() ident.Id {
	panic(fmt.Errorf("get state not supported for property %v", p.prop.GetId()))
}
func (p panicValue) SetState(ident.Id) error {
	panic(fmt.Errorf("set state not supported for property %v", p.prop.GetId()))
}
func (p panicValue) GetObject() ident.Id {
	panic(fmt.Errorf("get object not supported for property %v", p.prop.GetId()))
}
func (p panicValue) SetObject(ident.Id) error {
	panic(fmt.Errorf("set object not supported for property %v", p.prop.GetId()))
}
