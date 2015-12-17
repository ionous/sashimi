package metal

import (
	"encoding/json"
	"github.com/ionous/sashimi/util/ident"
	"io"
)

// ObjectValueMap provides a default implementation of ObjectValue
type ObjectValueMap map[string]interface{}

// GetValue succeeds if SetValue was called on a corresponding obj.field.
func (m ObjectValueMap) GetValue(obj, field ident.Id) (ret interface{}, okay bool) {
	n := obj.String() + "." + field.String()
	if value, ok := m[n]; ok {
		ret, okay = value, ok
	}
	return
}

// SetValue always succeeds, storing the passed value to the map at obj.field.
func (m ObjectValueMap) SetValue(obj, field ident.Id, value interface{}) (err error) {
	n := obj.String() + "." + field.String()
	m[n] = value
	return
}

func (m ObjectValueMap) Save(w io.Writer) (err error) {
	return json.NewEncoder(w).Encode(m)
}

func (m ObjectValueMap) Load(r io.Reader) (err error) {
	return json.NewDecoder(r).Decode(m)
}
