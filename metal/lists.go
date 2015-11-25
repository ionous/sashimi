package metal

import (
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
	"reflect"
)

// implements meta.Values
type arrayValues struct {
	*propBase
	// valueNum constructs a new meta.Value wrapper
	// (a num, a text, an object) for the passed index.
	valueNum func(int) meta.Value
}

func (ar arrayValues) NumValue() int {
	slice := reflect.ValueOf(ar.getGeneric())
	return slice.Len()
}

func (ar arrayValues) ValueNum(i int) meta.Value {
	return ar.valueNum(i)
}

func (ar arrayValues) ClearValues() error {
	empty := ar.mdl.getZero(ar.prop)
	return ar.set(empty)
}

func (ar arrayValues) AppendNum(v float32) error {
	slice := ar.getGeneric().([]float32)
	return ar.set(append(slice, v))
}

func (ar arrayValues) AppendText(v string) error {
	slice := ar.getGeneric().([]string)
	return ar.set(append(slice, v))
}

func (ar arrayValues) AppendObject(v ident.Id) error {
	slice := ar.getGeneric().([]ident.Id)
	return ar.set(append(slice, v))
}
