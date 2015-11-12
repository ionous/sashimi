package memory

import (
	"github.com/ionous/sashimi/util/ident"
)

// this currently assumes lists are empty by default
// we could add a return values from .get that says where the value came from
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
	val := el.get()
	return val.(float32)
}
func (el numElement) SetNum(f float32) error {
	slice := el.slice()
	slice[el.index] = f
	return el.set(slice)
}

func (el textElement) GetText() string {
	val := el.get()
	return val.(string)
}
func (el textElement) SetText(t string) error {
	slice := el.slice()
	slice[el.index] = t
	return el.set(slice)
}

func (el objectElement) GetObject() ident.Id {
	val := el.get()
	return val.(ident.Id)
}
func (el objectElement) SetObject(n ident.Id) error {
	slice := el.slice()
	slice[el.index] = n
	return el.set(slice)
}
func (el elementValue) slice() []interface{} {
	return el.propBase.get().([]interface{})
}

func (el elementValue) get() GenericValue {
	slice := el.slice()
	return slice[el.index]
}
