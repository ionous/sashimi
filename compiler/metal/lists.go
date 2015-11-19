package metal

import (
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
)

// implements meta.Values
type arrayValues struct {
	*propBase
	valueNum func(int) meta.Value
}

func (ar arrayValues) NumValue() int {
	slice := ar.get().([]interface{})
	return len(slice)
}

func (ar arrayValues) ValueNum(i int) meta.Value {
	return ar.valueNum(i)
}

func (ar arrayValues) ClearValues() error {
	return ar.set([]interface{}{})
}

func (ar arrayValues) AppendNum(f float32) error {
	return ar.append(f)
}

func (ar arrayValues) AppendText(t string) error {
	return ar.append(t)
}

func (ar arrayValues) AppendObject(n ident.Id) error {
	return ar.append(n)
}

func (ar arrayValues) append(x interface{}) error {
	slice := ar.get().([]interface{})
	slice = append(slice, GenericValue(x))
	return ar.set(slice)
}
