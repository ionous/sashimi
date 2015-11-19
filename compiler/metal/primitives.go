package metal

import "github.com/ionous/sashimi/util/ident"

type numValue struct{ panicValue }

func (p numValue) SetNum(f float32) error {
	return p.set(f)
}
func (p numValue) GetNum() float32 {
	return p.get().(float32)
}

type numValues struct{ panicValue }

func (p numValues) SetNum(f float32) error {
	return p.set(f)
}
func (p numValues) GetNum() float32 {
	return p.get().(float32)
}

type textValue struct{ panicValue }

func (p textValue) GetText() string {
	return p.get().(string)
}
func (p textValue) SetText(t string) error {
	return p.set(t)
}

type pointerValue struct {
	panicValue
}

func (p pointerValue) GetObject() ident.Id {
	return p.get().(ident.Id)
}
func (p pointerValue) SetObject(n ident.Id) error {
	return p.set(n)
}
