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
	panic("use list eval")
}
func (el *numElement) SetNum(v float64) error {
	panic("use list eval")
}

func (el *textElement) GetText() string {
	panic("use list eval")
}
func (el *textElement) SetText(v string) error {
	panic("use list eval")
}

func (el *objectElement) GetObject() ident.Id {
	panic("use list eval")
}
func (el *objectElement) SetObject(v ident.Id) error {
	panic("use list eval")
}
