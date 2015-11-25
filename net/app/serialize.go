package app

import (
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/net/resource"
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
func (s *ObjectSerializer) IsKnown(gobj meta.Instance) bool {
	return gobj != nil && s.known[gobj.GetId()]
}

//
// Add the object to the passed document data as the primary object.
//
// NOTE: unlike, jsonapi/rest we omit the existance and contents of relationships.
// The client will ask explicitly ask for the relation information it wants:
// ex. /games/{session}/actors/player/inventory
//
func (s *ObjectSerializer) SerializeObject(out resource.IBuildObjects, gobj meta.Instance, force bool) (obj *resource.Object) {
	if s.known.SetKnown(gobj) || force {
		obj = s.NewObject(out, gobj)
		//
		states := []string{}
		for i := 0; i < gobj.NumProperty(); i++ {
			prop := gobj.PropertyNum(i)
			pid := jsonId(prop.GetId())
			switch prop.GetType() {
			case meta.NumProperty:
				obj.SetAttr(pid, prop.GetValue().GetNum())
			case meta.TextProperty:
				obj.SetAttr(pid, prop.GetValue().GetText())
			case meta.StateProperty:
				choice := jsonId(prop.GetValue().GetState())
				states = append(states, choice)
				obj.SetAttr(pid, choice)
			case meta.ObjectProperty:
				obj.SetAttr(pid, jsonId(prop.GetValue().GetObject()))
			default:
				// ignore arrays for now....
			}
		}

		obj.SetMeta("name", gobj.GetOriginalName())
		obj.SetMeta("states", states)
	}
	return obj
}

//
// Add a reference to the passed object into the passed refs list,
// with a full seriaization into includes if the object is newly known.
//
func (s *ObjectSerializer) AddObjectRef(out resource.IBuildObjects, gobj meta.Instance, include resource.IBuildObjects) (obj *resource.Object) {
	s.SerializeObject(include, gobj, false)
	return s.NewObject(out, gobj)
}

func (s *ObjectSerializer) NewObject(out resource.IBuildObjects, gobj meta.Instance) (obj *resource.Object) {
	return out.NewObject(jsonId(gobj.GetId()), jsonId(gobj.GetParentClass().GetId()))
}
