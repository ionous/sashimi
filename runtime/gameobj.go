package runtime

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/model/table"
	"github.com/ionous/sashimi/util/ident"
	"reflect"
)

// GameObject
type GameObject struct {
	id     ident.Id       // unique id, matches instance info's ids.
	cls    *M.ClassInfo   // for property set, etc. access
	vals   TemplateValues // runtime gobj are key'd by string for go's templates
	temps  TemplatePool   // FIX? cache for templates.... probably should nix this.
	tables table.Tables
}

//
// GameObjects maps model instance id to runtime game object class.
//
type GameObjects map[ident.Id]*GameObject

//
// Id uniquely identifies this object.
//
func (gobj *GameObject) Id() ident.Id {
	return gobj.id
}

//
// Class of this game object.
//
func (gobj *GameObject) Class() *M.ClassInfo {
	return gobj.cls
}

//
// String representation of the object's id.
//
func (gobj *GameObject) String() string {
	return gobj.id.String()
}

// nil if it didnt exist, which beacuse the gobj for the instances are "flattened"
// and because nil isn't used for the default value of anything, should be a fine signal.
func (gobj *GameObject) Value(id ident.Id) interface{} {
	return gobj.vals[id.String()]
}

//
// set, but only if type of the current value at name matches the passed value
//
func (gobj *GameObject) SetValue(id ident.Id, val interface{}) (old interface{}, okay bool) {
	if v, had := gobj.vals[id.String()]; had &&
		reflect.TypeOf(v) == reflect.TypeOf(val) {
		gobj.setDirect(id, val)
		old = had
		okay = true
	}
	return old, okay
}

//
func (gobj *GameObject) setDirect(id ident.Id, value interface{}) {
	gobj.vals[id.String()] = value
}

//
func (gobj *GameObject) removeDirect(id ident.Id) {
	delete(gobj.vals, id.String())
}

// set the property value.
func (gobj *GameObject) setValue(prop M.IProperty, val interface{}) (err error) {
	switch prop := prop.(type) {
	case M.EnumProperty:
		if choice, e := prop.IndexToChoice(val.(int)); e != nil {
			err = e
		} else {
			gobj.setDirect(prop.GetId(), choice)
			gobj.setDirect(choice, true)
		}

	case M.NumProperty:
		gobj.setDirect(prop.GetId(), val)

	case M.PointerProperty:
		gobj.setDirect(prop.GetId(), val)

	case M.RelativeProperty:
		if table, ok := gobj.tables[prop.Relation]; !ok {
			err = fmt.Errorf("couldn't find table", prop.Relation)
		} else {
			rel := RelativeValue{gobj.Id(), prop, table}
			gobj.setDirect(prop.GetId(), rel)
		}

	case M.TextProperty:
		gobj.setDirect(prop.GetId(), val)
		// TBD: when to parse this? maybe options? here for early errors.
		str := val.(string)
		if e := gobj.temps.New(prop.GetId().String(), str); e != nil {
			err = e
		}

	default:
		err = fmt.Errorf("internal error: unknown property type %s:%T", prop, prop)
	}
	return err
}
