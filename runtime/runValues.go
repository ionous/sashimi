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
	key := id.String()
	this.data[key] = value
}

//
// func (this RuntimeValues) Any(id M.StringId) (interface{}, bool) {
// 	return this.data[id.String()]
// }

//
func (this RuntimeValues) Num(id M.StringId) (float32, bool) {
	value, ok := this.data[id.String()].(float32)
	return value, ok
}

//
func (this RuntimeValues) Text(id M.StringId) (string, bool) {
	value, ok := this.data[id.String()].(string)
	return value, ok
}

//
// pass the property id: ex. ScoredProperty to determine the current selection
//
func (this RuntimeValues) Choice(id M.StringId) (M.StringId, bool) {
	value, ok := this.data[id.String()].(M.StringId)
	return value, ok
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
