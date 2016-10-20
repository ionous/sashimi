package metal

import (
	"github.com/ionous/sashimi/util/ident"
)

// this currently assumes lists are empty by default
// we could add a return values from .getIdx that says where the value came from
// and clone the array if it came from the instance.
type elementValue struct {
	PanicValue
	index int
}

type numElement struct {
	*elementValue
}
type textElement struct {
	*elementValue
}
type objectElement struct {
	*elementValue
}

func (el *numElement) GetNum() float64 {
	slice := el.GetGeneric().([]float64)
	return slice[el.index]
}
func (el *numElement) SetNum(v float64) error {
	slice := el.GetGeneric().([]float64)
	slice[el.index] = v
	return el.SetGeneric(slice)
}

func (el *textElement) GetText() string {
	slice := el.GetGeneric().([]string)
	return slice[el.index]
}
func (el *textElement) SetText(v string) error {
	slice := el.GetGeneric().([]string)
	slice[el.index] = v
	return el.SetGeneric(slice)
}

func (el *objectElement) GetObject() ident.Id {
	slice := el.GetGeneric().([]ident.Id)
	return slice[el.index]
}
func (el *objectElement) SetObject(v ident.Id) error {
	slice := el.GetGeneric().([]ident.Id)
	slice[el.index] = v
	return el.SetGeneric(slice)
}
