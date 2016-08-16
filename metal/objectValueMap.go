package metal

import "github.com/ionous/sashimi/util/ident"

// ObjectValueMap provides a default implementation of ObjectValue
type ObjectValueMap map[string]interface{}

// GetValue succeeds if SetValue was called on a corresponding obj.field.
func (m ObjectValueMap) GetValue(obj, field ident.Id) (interface{}, bool) {
	value, ok := m[string(obj)+"."+string(field)]
	return value, ok
}

// SetValue always succeeds, storing the passed value to the map at obj.field.
func (m ObjectValueMap) SetValue(obj, field ident.Id, value interface{}) (err error) {
	m[string(obj)+"."+string(field)] = value
	return
}
