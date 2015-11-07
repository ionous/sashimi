package app

import (
	"github.com/ionous/sashimi/net/resource"
	R "github.com/ionous/sashimi/runtime"
	"github.com/ionous/sashimi/runtime/api"
)

//
// Generates object documents which include a turn counter.
// A document contains refs to objects, and optionally includes serialization of those referenced objects.
// This uses state to determine which objects we've already told the client about.
//
type ObjectSerializer struct {
	known KnownObjects
}

//
func NewObjectSerializer() *ObjectSerializer {
	return &ObjectSerializer{known: make(KnownObjects)}
}

//
func (s *ObjectSerializer) IsKnown(gobj *R.GameObject) bool {
	return gobj != nil && s.known[gobj.Id()]
}

//
// Add the object to the passed document data as the primary object.
//
// NOTE: unlike, jsonapi/rest we omit the existance and contents of relationships.
// The client will ask explicitly ask for the relation information it wants:
// ex. /games/{session}/actors/player/inventory
//
func (s *ObjectSerializer) SerializeObject(out resource.IBuildObjects, gobj *R.GameObject, force bool) (obj *resource.Object) {
	if s.known.SetKnown(gobj) || force {
		obj = s.NewObject(out, gobj)
		//
		states := []string{}
		// OBJECT FIX: should be from objcet
		cls := gobj.Class()
		for i := 0; i < cls.NumProperty(); i++ {
			prop := cls.PropertyNum(i)
			if val := gobj.Value(prop.GetId()); val != nil {
				id := jsonId(prop.GetId())
				switch prop.GetType() {
				case api.NumProperty:
					obj.SetAttr(id, prop.GetValue().GetNum())
				case api.TextProperty:
					obj.SetAttr(id, prop.GetValue().GetText())
				case api.StateProperty:
					obj.SetAttr(id, jsonId(prop.GetValue().GetState()))
				case api.ObjectProperty:
					obj.SetAttr(id, jsonId(prop.GetValue().GetObject()))
				default:
					// ignore arrays for now....
				}
			}
		}
		obj.
			SetMeta("name", gobj.String()).
			SetMeta("states", states)
	}
	return obj
}

//
// Add a reference to the passed object into the passed refs list,
// with a full seriaization into includes if the object is newly known.
//
func (s *ObjectSerializer) AddObjectRef(out resource.IBuildObjects, gobj *R.GameObject, include resource.IBuildObjects) (obj *resource.Object) {
	s.SerializeObject(include, gobj, false)
	return s.NewObject(out, gobj)
}

func (s *ObjectSerializer) NewObject(out resource.IBuildObjects, gobj *R.GameObject) (obj *resource.Object) {
	return out.NewObject(jsonId(gobj.Id()), jsonId(gobj.Class().GetId()))
}
