package app

import (
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/net/resource"
	"github.com/ionous/sashimi/util/ident"
)

// ObjectSerializer generates object documents which include a turn counter.
// A document contains refs to objects, and optionally includes serialization of those referenced objects.
type ObjectSerializer struct {
	KnownObjects
	includes Includes
}

// Includes all objects referenced by the CommandOutput.
type Includes map[ident.Id]meta.Instance

func (inc Includes) Include(gobj meta.Instance) {
	inc[gobj.GetId()] = gobj
}

// NewObjectSerializer uses known objects to determine which objects we've already told the client about.
func NewObjectSerializer(known KnownObjects) *ObjectSerializer {
	return &ObjectSerializer{known, make(Includes)}
}

// TryObjectRef only creates an object ref if the object is already known.
func (s *ObjectSerializer) TryObjectRef(gobj meta.Instance) (ret *resource.Object, okay bool) {
	if s.IsKnown(gobj) {
		ret = s.NewObjectRef(gobj)
		okay = true
	}
	return
}

// NewObjectRef always adds the object to the includes
func (s *ObjectSerializer) NewObjectRef(gobj meta.Instance) *resource.Object {
	s.includes.Include(gobj)
	out := resource.ObjectList{}
	return out.NewObject(jsonId(gobj.GetId()), jsonId(gobj.GetParentClass().GetId()))
}

// Flush returns the objects we needed to include.
func (s *ObjectSerializer) Flush() Includes {
	ret := s.includes
	s.includes = make(Includes)
	return ret
}

// SerializeObject to the passed document data as the primary object.
// NOTE: unlike, jsonapi/rest we omit the existance and contents of relationships.
// The client explicitly asks for the relation information it wants:
// ex. /games/{session}/actors/player/inventory
func (s *ObjectSerializer) SerializeObject(out resource.IBuildObjects, gobj meta.Instance, force bool) (obj *resource.Object) {
	if s.SetKnown(gobj) || force {
		obj = out.NewObject(jsonId(gobj.GetId()), jsonId(gobj.GetParentClass().GetId()))
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

// AddObjectRef adds the passed object into the passed list of references,
// with a full seriaization into includes if the object is newly known.
func (s *ObjectSerializer) AddObjectRef(out resource.IBuildObjects, gobj meta.Instance, include resource.IBuildObjects) (obj *resource.Object) {
	s.SerializeObject(include, gobj, false)
	return out.NewObject(jsonId(gobj.GetId()), jsonId(gobj.GetParentClass().GetId()))
}

// NewObject adds a reference to the object into the current document.
// func (s *ObjectSerializer) NewObject(out resource.IBuildObjects, gobj meta.Instance) (obj *resource.Object) {
// 	return out.NewObject(jsonId(gobj.GetId()), jsonId(gobj.GetParentClass().GetId()))
//}
