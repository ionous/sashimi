package app

import (
	"github.com/ionous/sashimi/net/resource"
	R "github.com/ionous/sashimi/runtime"
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
	return this.known[gobj.Id()]
}

//
// Add the object to the passed document data as the primary object.
//
// NOTE: unlike, jsonapi/rest we omit the existance and contents of relationships.
// The client will ask explicitly ask for the relation information it wants:
// ex. /games/{session}/actors/player/inventory
//
func (this *ObjectSerializer) SerializeObject(out resource.IBuildObjects, gobj *R.GameObject, force bool) (obj *resource.Object) {
	if this.known.SetKnown(gobj.Id()) || force {
		obj = this.NewObject(out, gobj)
		//
		states := []string{}
		for prop, _ := range gobj.Class().AllProperties() {
			// FIX: this shouldnt require three map lookups
			if choice, ok := gobj.Choice(prop); ok {
				states = append(states, jsonId(choice))
			} else if text, ok := gobj.Text(prop); ok {
				obj.SetAttr(jsonId(prop), text)
			} else if num, ok := gobj.Num(prop); ok {
				obj.SetAttr(jsonId(prop), num)
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
	return out.NewObject(jsonId(gobj.Id()), jsonId(gobj.Class().Id()))
}

// //
// // internal helper for iterating over class hierarchy.
// //
// func classList(cls *M.ClassInfo) (ret []*M.ClassInfo) {
// 	ret = append(ret, cls)
// 	if par := cls.Parent(); par != nil {
// 		ret = append(ret, classList(par)...)
// 	}
// 	return ret
// }