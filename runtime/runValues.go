package runtime

import (
	M "github.com/ionous/sashimi/model"
	"reflect"
)

type RuntimeValues struct {
	data TemplateValues // runtime values are key'd by string for go's templates
}

//
func NewRuntimeValues() RuntimeValues {
	return RuntimeValues{make(map[string]interface{})}
}

//
func (this RuntimeValues) removeDirect(id M.StringId) {
	delete(this.data, id.String())
}

//
func (this RuntimeValues) setDirect(id M.StringId, value interface{}) {
	this.data[id.String()] = value
}

// nil if it didnt exist, which beacuse the values for the instances are "flattened"
// and because nil isn't used for the default value of anything, should be a fine signal.
func (this RuntimeValues) GetValue(id M.StringId) interface{} {
	return this.data[id.String()]
}

//
// set, but only if type of the current value at name matches the passed value
//
func (this RuntimeValues) SetValue(id M.StringId, val interface{}) (old interface{}, okay bool) {
	if v, had := this.data[id.String()]; had &&
		reflect.TypeOf(v) == reflect.TypeOf(val) {
		this.setDirect(id, val)
		old = had
		okay = true
	}
	return old, okay
}
