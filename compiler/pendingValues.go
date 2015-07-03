package compiler

import (
	M "github.com/ionous/sashimi/model"
	"reflect" // for type checking
)

//
// Property values for a partial instance.
//
type PendingValues map[M.StringId]M.Variant

//
// lockSet sets the passed property to value, erroring if the there is a type mismatch or a value conflict.
// subsequent calls to the same property must share the saee value.
//
func (vals PendingValues) lockSet(inst, prop M.StringId, nilVal, want interface{}) (err error) {
	if reflect.TypeOf(nilVal) != reflect.TypeOf(want) {
		err = SetValueMismatch(inst, prop, nilVal, want)
	} else if curr, have := vals[prop]; have && curr != want {
		err = SetValueChanged(inst, prop, curr, want)
	} else {
		vals[prop] = want
	}
	return err
}
