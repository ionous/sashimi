package metal

import (
	"fmt"
	M "github.com/ionous/sashimi/compiler/model"
	"github.com/ionous/sashimi/compiler/model/table"
	"github.com/ionous/sashimi/meta/tests"
	"github.com/ionous/sashimi/util/ident"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// go test -run TestMemory
func TestMemory(t *testing.T) {
	v := make(ObjectValueMap)
	m := &M.Model{
		Classes:      make(M.Classes),
		Instances:    make(M.Instances),
		Enumerations: make(M.Enumerations),
	}
	numProp, numsProp := ident.MakeId("Num"), ident.MakeId("Nums")
	textProp, textsProp := ident.MakeId("Text"), ident.MakeId("Texts")
	stateProp := ident.MakeId("State")
	objectProp, objectsProp := ident.MakeId("Object"), ident.MakeId("Objects")
	m.Enumerations[stateProp] = &M.EnumModel{[]ident.Id{ident.MakeId("No"), ident.MakeId("Yes")}}

	makeClass := func(single, plural string) *M.ClassModel {
		clsId := ident.MakeId(plural)
		//
		c := &M.ClassModel{
			Id:       clsId,
			Plural:   plural,
			Singular: single,
			Properties: []M.PropertyModel{
				M.PropertyModel{Id: numProp, Name: "Num", Type: M.NumProperty},
				M.PropertyModel{Id: textProp, Name: "Text", Type: M.TextProperty},
				M.PropertyModel{Id: objectProp, Name: "Object", Type: M.PointerProperty, Relates: clsId},
				//
				M.PropertyModel{Id: numsProp, Name: "Nums", Type: M.NumProperty, IsMany: true},
				M.PropertyModel{Id: textsProp, Name: "Texts", Type: M.TextProperty, IsMany: true},
				M.PropertyModel{Id: objectsProp, Name: "Objects", Type: M.PointerProperty, Relates: clsId, IsMany: true},
				//
				M.PropertyModel{Id: stateProp, Name: "State", Type: M.EnumProperty},
			},
			// Constraints,
		}
		m.Classes[clsId] = c
		return c
	}
	makeInstance := func(name string, c *M.ClassModel) *M.InstanceModel {
		inst := &M.InstanceModel{
			Id:    ident.MakeId(name),
			Class: c.Id,
			Name:  name,
			// Values
		}
		m.Instances[inst.Id] = inst
		return inst
	}
	i := makeInstance("i", makeClass("test class", "test classes"))
	makeInstance("x", makeClass("test other", "test others"))
	//m.PrintModel(t.Log)
	//t.Fatal("")
	mem := NewMetal(m, v, make(table.Tables))
	// test plural
	if cls, ok := mem.GetClass(ident.MakeId("TestClasses")); assert.True(t, ok, "get test class") {
		if p, ok := cls.FindProperty("plural"); assert.True(t, ok) {
			require.EqualValues(t, "test classes", p.GetValue().GetText())
		}
		if p, ok := cls.FindProperty("singular"); assert.True(t, ok) {
			require.EqualValues(t, "test class", p.GetValue().GetText())
		}
	}
	// test the api
	tests.ApiTest(t, mem, i.Id)
	// test that things really changed
	if res, ok := v.GetValue(i.Id, numProp); !assert.True(t, ok, fmt.Sprintf("missing (%s.%s)", i.Id, numProp)) {
		t.Log(v)
		t.Fatal()
	} else {
		require.EqualValues(t, float32(32), res)
	}
	if res, ok := v.GetValue(i.Id, textProp); assert.True(t, ok) {
		require.EqualValues(t, "text", res)
	}
	if res, ok := v.GetValue(i.Id, stateProp); assert.True(t, ok) {
		// 1= No, 2= Yes.
		// Why *does** GetValue return an integer anyway?
		require.EqualValues(t, 2, res)

	}
	if res, ok := v.GetValue(i.Id, objectProp); assert.True(t, ok) {
		require.EqualValues(t, i.Id, res)
	}
	//
	if res, ok := v.GetValue(i.Id, numsProp); assert.True(t, ok) {
		require.Contains(t, res, float32(32))
	}
	if res, ok := v.GetValue(i.Id, textsProp); assert.True(t, ok) {
		require.Contains(t, res, "text")
	}
	if res, ok := v.GetValue(i.Id, objectProp); assert.True(t, ok) {
		require.Contains(t, res, i.Id)
	}
}
