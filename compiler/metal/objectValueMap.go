package metal

import (
	"github.com/ionous/sashimi/util/ident"
)

// ObjectValueMap provides a default implementation of ObjectValue
type ObjectValueMap map[string]GenericValue

// GetValue succeeds if SetValue was called on a corresponding obj.field.
func (m ObjectValueMap) GetValue(obj, field ident.Id) (ret GenericValue, okay bool) {
	n := obj.String() + "." + field.String()
	if value, ok := m[n]; ok {
		ret, okay = value, ok
	}
	return
}

// SetValue always succeeds, storing the passed value to the map at obj.field.
func (m ObjectValueMap) SetValue(obj, field ident.Id, value GenericValue) (err error) {
	n := obj.String() + "." + field.String()
	m[n] = value
	return
}
