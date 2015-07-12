package runtime

import (
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/util/ident"
	"reflect"
)

type RuntimeValues struct {
	data   TemplateValues // runtime values are key'd by string for go's templates
	temps  TemplatePool
	tables M.TableRelations
}

//
func NewRuntimeValues(tables M.TableRelations) RuntimeValues {
	return RuntimeValues{make(map[string]interface{}), make(TemplatePool), tables}
}

//
func (this RuntimeValues) removeDirect(id ident.Id) {
	delete(this.data, id.String())
}

//
func (this RuntimeValues) setDirect(id ident.Id, value interface{}) {
	this.data[id.String()] = value
}

// nil if it didnt exist, which beacuse the values for the instances are "flattened"
// and because nil isn't used for the default value of anything, should be a fine signal.
func (this RuntimeValues) GetValue(id ident.Id) interface{} {
	return this.data[id.String()]
}

//
// set, but only if type of the current value at name matches the passed value
//
func (this RuntimeValues) SetValue(id ident.Id, val interface{}) (old interface{}, okay bool) {
	if v, had := this.data[id.String()]; had &&
		reflect.TypeOf(v) == reflect.TypeOf(val) {
		this.setDirect(id, val)
		old = had
		okay = true
	}
	return old, okay
}
