package metal

import "github.com/ionous/sashimi/util/ident"

// ObjectValueMap provides a default implementation of ObjectValue
type ObjectValueMap map[string]interface{}

// GetValue succeeds if SetValue was called on a corresponding obj.field.
func (m ObjectValueMap) GetValue(obj, field ident.Id) (ret interface{}, okay bool) {
	n := string(obj) + "." + string(field)
	if value, ok := m[n]; ok {
		ret, okay = value, ok
	}
	return
}

// SetValue always succeeds, storing the passed value to the map at obj.field.
func (m ObjectValueMap) SetValue(obj, field ident.Id, value interface{}) (err error) {
	n := string(obj) + "." + string(field)
	m[n] = value
	return
}
