package internal

import (
	M "github.com/ionous/sashimi/compiler/xmodel"
	"github.com/ionous/sashimi/util/ident"
	"reflect" // for type checking
)

// Property values for a partial instance.
type PendingValues map[ident.Id]M.Variant

// lockSet sets the passed property to value, erroring if the there is a type mismatch or a value conflict.
// subsequent calls to the same property must share the same value.
func (vals PendingValues) lockSet(inst, prop ident.Id, iface, want interface{}) (err error) {
	if !reflect.TypeOf(want).Implements(reflect.TypeOf(iface).Elem()) {
		err = SetValueMismatch("lockSet", inst, prop, iface, want)
	} else if curr, have := vals[prop]; have && curr != want {
		err = SetValueChanged(inst, prop, curr, want)
	} else {
		vals[prop] = want
	}
	return err
}
