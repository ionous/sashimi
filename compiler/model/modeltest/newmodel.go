package modeltest

import (
	M "github.com/ionous/sashimi/compiler/model"
	"github.com/ionous/sashimi/util/ident"
)

var TestInstance = ident.MakeId("i")
var OtherInstance = ident.MakeId("x")
var NumProp, NumsProp = ident.MakeId("Num"), ident.MakeId("Nums")
var TextProp, TextsProp = ident.MakeId("Text"), ident.MakeId("Texts")
var StateProp = ident.MakeId("State")
var ObjectProp, ObjectsProp = ident.MakeId("Object"), ident.MakeId("Objects")

// NewModel creates a compiler model for simple testing purposes.
func NewModel() *M.Model {
	m := &M.Model{
		Classes:      make(M.Classes),
		Instances:    make(M.Instances),
		Enumerations: make(M.Enumerations),
	}
	m.Enumerations[StateProp] = &M.EnumModel{[]ident.Id{ident.MakeId("No"), ident.MakeId("Yes")}}
	makeClass := func(single, plural string) *M.ClassModel {
		clsId := ident.MakeId(plural)
		//
		c := &M.ClassModel{
			Id:       clsId,
			Plural:   plural,
			Singular: single,
			Properties: []M.PropertyModel{
				M.PropertyModel{Id: NumProp, Name: "Num", Type: M.NumProperty},
				M.PropertyModel{Id: TextProp, Name: "Text", Type: M.TextProperty},
				M.PropertyModel{Id: ObjectProp, Name: "Object", Type: M.PointerProperty, Relates: clsId},
				//
				M.PropertyModel{Id: NumsProp, Name: "Nums", Type: M.NumProperty, IsMany: true},
				M.PropertyModel{Id: TextsProp, Name: "Texts", Type: M.TextProperty, IsMany: true},
				M.PropertyModel{Id: ObjectsProp, Name: "Objects", Type: M.PointerProperty, Relates: clsId, IsMany: true},
				//
				M.PropertyModel{Id: StateProp, Name: "State", Type: M.EnumProperty},
			},
			// Constraints,
		}
		m.Classes[clsId] = c
		return c
	}
	makeInstance := func(id ident.Id, c *M.ClassModel) *M.InstanceModel {
		inst := &M.InstanceModel{
			Id:    id,
			Class: c.Id,
			Name:  id.String(),
			// Values
		}
		m.Instances[inst.Id] = inst
		return inst
	}
	makeInstance(TestInstance, makeClass("test class", "test classes"))
	makeInstance(OtherInstance, makeClass("test other", "test others"))
	return m
}
