package runtime

import (
	"fmt"
	"github.com/ionous/sashimi/model/table"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/util/errutil"
	"github.com/ionous/sashimi/util/ident"
	"reflect"
)

// GameObject
type GameObject struct {
	id     ident.Id // not all game objects come from instances
	cls    api.Class
	vals   RuntimeValues
	tables table.Tables
}

func NewGameObject(
	mdl api.Model,
	objId ident.Id,
	cls api.Class,
	proto api.Prototype,
	tables table.Tables,
) (_ *GameObject, err error,
) {
	gobj := &GameObject{objId, cls, make(RuntimeValues), tables}
	// store-fix: with fallback, is this needed?
	for i := 0; i < proto.NumProperty(); i++ {
		prop := proto.PropertyNum(i)
		switch prop.GetType() {
		case api.StateProperty:
			choice := prop.GetValue().GetState()
			gobj.setDirect(prop.GetId(), choice)
			gobj.setDirect(choice, true)

		case api.NumProperty:
			gobj.setDirect(prop.GetId(), prop.GetValue().GetNum())
		case api.TextProperty:
			gobj.setDirect(prop.GetId(), prop.GetValue().GetText())
		case api.ObjectProperty:
			gobj.setDirect(prop.GetId(), prop.GetValue().GetObject())
		case api.ObjectProperty | api.ArrayProperty:
			// ignore
		default:
			err = errutil.Append(err, fmt.Errorf("unknown property type %s:%v", prop, prop.GetType()))
		}
	}
	return gobj, err
}

type RuntimeValues map[string]interface{}

// GameObjects maps model instance id to runtime game object class.
type GameObjects map[ident.Id]*GameObject

// Id uniquely identifies this object.
func (gobj *GameObject) Id() ident.Id {
	return gobj.id
}

// Class of this game object.
func (gobj *GameObject) Class() api.Class {
	return gobj.cls
}

// String representation of the object's id.
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
		old, okay = v, true
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
