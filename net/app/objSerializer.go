package app

import (
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/net/resource"
	"github.com/ionous/sashimi/standard/framework" //containment includes
	"github.com/ionous/sashimi/util/ident"
)

// ObjSerializer generates object documents which include a turn counter.
// A document contains refs to out, and optionally out serialization of those referenced out.
type ObjSerializer struct {
	mdl   meta.Model
	out   resource.ObjectList
	known map[ident.Id]bool
}

func NewObjSerializer(m meta.Model, out resource.ObjectList) *ObjSerializer {
	return &ObjSerializer{m, out, make(map[ident.Id]bool)}
}

// Include serializes the object and its contents
func (o *ObjSerializer) Include(gobj meta.Instance) {
	if o.SerializeObject(gobj) {
		for _, rel := range framework.Containment {
			if prop, ok := gobj.GetProperty(rel); ok {
				values := prop.GetValues()
				for i := 0; i < values.NumValue(); i++ {
					id := values.ValueNum(i).GetObject()
					if inst, ok := o.mdl.GetInstance(id); ok {
						o.Include(inst)
					}
				}
			}
		}
	}
}

// SerializeObject adds the object to the wraped resource ( if it hasnt been done before )
func (o *ObjSerializer) SerializeObject(gobj meta.Instance) (added bool) {
	if id := gobj.GetId(); !o.known[id] {
		o.known[id] = true
		o.out.AddObject(SerializeObject(gobj))
		added = true
	}
	return
}

// SerializeObject to the passed resource.
// NOTE: unlike, jsonapi/rest we omit relationships and arrays.
func SerializeObject(gobj meta.Instance) *resource.Object {
	obj := NewObjectRef(gobj)
	//
	states := []string{}
	for i := 0; i < gobj.NumProperty(); i++ {
		prop := gobj.PropertyNum(i)
		pid := jsonId(prop.GetId())
		switch t := prop.GetType(); t {
		case meta.NumProperty:
			v := prop.GetValue() // note, sharing this "GetValue" panics, because not all types support GetValue
			obj.SetAttr(pid, v.GetNumber())
		case meta.TextProperty:
			v := prop.GetValue()
			obj.SetAttr(pid, v.GetText())
		case meta.StateProperty:
			v := prop.GetValue()
			choice := jsonId(v.GetState())
			states = append(states, choice)
			obj.SetAttr(pid, choice)
		case meta.ObjectProperty:
			v := prop.GetValue()
			id := jsonId(v.GetObject())
			obj.SetAttr(pid, id)
		default:
			// ignore arrays for now....
		}
	}
	obj.SetMeta("name", gobj.GetOriginalName())
	obj.SetMeta("states", states)
	return obj
}
