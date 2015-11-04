package app

import (
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/net/resource"
	R "github.com/ionous/sashimi/runtime"
	"github.com/ionous/sashimi/util/ident"
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
func (this *ObjectSerializer) IsKnown(gobj *R.GameObject) bool {
	return gobj != nil && this.known[gobj.Id()]
}

//
// Add the object to the passed document data as the primary object.
//
// NOTE: unlike, jsonapi/rest we omit the existance and contents of relationships.
// The client will ask explicitly ask for the relation information it wants:
// ex. /games/{session}/actors/player/inventory
//
func (this *ObjectSerializer) SerializeObject(out resource.IBuildObjects, gobj *R.GameObject, force bool) (obj *resource.Object) {
	if this.known.SetKnown(gobj) || force {
		obj = this.NewObject(out, gobj)
		//
		states := []string{}
		for propId, prop := range gobj.Class().AllProperties() {
			if val := gobj.Value(propId); val != nil {
				switch prop.(type) {
				case M.EnumProperty:
					choice := val.(ident.Id)
					states = append(states, jsonId(choice))
				default:
					obj.SetAttr(jsonId(propId), val)
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
func (this *ObjectSerializer) AddObjectRef(out resource.IBuildObjects, gobj *R.GameObject, include resource.IBuildObjects) (obj *resource.Object) {
	this.SerializeObject(include, gobj, false)
	return this.NewObject(out, gobj)
}

func (this *ObjectSerializer) NewObject(out resource.IBuildObjects, gobj *R.GameObject) (obj *resource.Object) {
	return out.NewObject(jsonId(gobj.Id()), jsonId(gobj.Class().Id))
}
