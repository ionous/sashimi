package metal

import (
	"github.com/ionous/sashimi/util/ident"
)

// this currently assumes lists are empty by default
// we could add a return values from .getIdx that says where the value came from
// and clone the array if it came from the instance.
type elementValue struct {
	panicValue
	index int
}

type numElement struct {
	elementValue
}
type textElement struct {
	elementValue
}
type objectElement struct {
	elementValue
}

func (el numElement) GetNum() float32 {
	slice := el.get().([]float32)
	return slice[el.index]
}
func (el numElement) SetNum(v float32) error {
	slice := el.get().([]float32)
	slice[el.index] = v
	return el.set(slice)
}

func (el textElement) GetText() string {
	slice := el.get().([]string)
	return slice[el.index]
}
func (el textElement) SetText(v string) error {
	slice := el.get().([]string)
	slice[el.index] = v
	return el.set(slice)
}

func (el objectElement) GetObject() ident.Id {
	slice := el.get().([]ident.Id)
	return slice[el.index]
}
func (el objectElement) SetObject(v ident.Id) error {
	slice := el.get().([]ident.Id)
	slice[el.index] = v
	return el.set(slice)
}
