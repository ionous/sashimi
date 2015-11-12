package memory

import (
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/util/ident"
)

// implements api.Values
type arrayValues struct {
	*propBase
	valueNum func(int) api.Value
}

func (ar arrayValues) NumValue() int {
	slice := ar.get().([]interface{})
	return len(slice)
}

func (ar arrayValues) ValueNum(i int) api.Value {
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
