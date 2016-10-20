package metal

import (
	"github.com/ionous/sashimi/util/errutil"
	"github.com/ionous/sashimi/util/ident"
)

// PanicValue implements the Value(s) interface:
// panicing on every get() and set(), and then
// specific property types override the specific methods they need:
// .text for text, num for num, etc.
type PanicValue struct {
	*propBase
}

func (p PanicValue) GetNum() float64 {
	panic(errutil.New("get num not supported for property", p.prop))
}
func (p PanicValue) SetNum(float64) error {
	panic(errutil.New("set num not supported for property", p.prop))
}
func (p PanicValue) AppendNum(float64) error {
	panic(errutil.New("append num not supported for property", p.prop))
}
func (p PanicValue) GetText() string {
	panic(errutil.New("get text not supported for property", p.prop))
}
func (p PanicValue) SetText(string) error {
	panic(errutil.New("set text not supported for property", p.prop))
}
func (p PanicValue) AppendText(string) error {
	panic(errutil.New("append text not supported for property", p.prop))
}
func (p PanicValue) GetObject() ident.Id {
	panic(errutil.New("get object not supported for property", p.prop))
}
func (p PanicValue) SetObject(ident.Id) error {
	panic(errutil.New("set object not supported for property", p.prop))
}
func (p PanicValue) SetRelation(ident.Id) error {
	panic(errutil.New("set relation not supported for property", p.prop))
}
func (p PanicValue) AppendObject(ident.Id) error {
	panic(errutil.New("append object not supported for property", p.prop))
}
func (p PanicValue) GetState() ident.Id {
	panic(errutil.New("get state not supported for property", p.prop))
}
func (p PanicValue) SetState(ident.Id) error {
	panic(errutil.New("set state not supported for property", p.prop))
}
