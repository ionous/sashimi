package memory

import (
	//"fmt"
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/model/table"
	"github.com/ionous/sashimi/runtime/api/tests"
	"github.com/ionous/sashimi/util/ident"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMemory(t *testing.T) {
	v := make(ObjectValueMap)
	m := &M.Model{
		Classes:   make(M.ClassMap),
		Instances: make(M.InstanceMap),
	}
	numProp := ident.MakeId("Num")
	textProp := ident.MakeId("Text")
	stateProp := ident.MakeId("State")
	objectProp := ident.MakeId("Object")
	enum := M.NewEnumeration([]string{"no", "yes"})

	makeClass := func(single, plural string) *M.ClassInfo {
		clsId := ident.MakeId(plural)
		c := &M.ClassInfo{
			Id:       clsId,
			Plural:   plural,
			Singular: single,
			Properties: M.PropertySet{
				"Num":   M.NumProperty{Id: numProp},
				"Text":  M.TextProperty{Id: textProp},
				"State": M.EnumProperty{Id: stateProp, Enumeration: enum},
				"Object": M.PointerProperty{Id: objectProp,
					Class: clsId},
			},
			// Constraints,
		}
		m.Classes[clsId] = c
		return c
	}
	makeInstance := func(name string, c *M.ClassInfo) *M.InstanceInfo {
		inst := &M.InstanceInfo{
			Id:    ident.MakeId(name),
			Class: c,
			Name:  name,
			// Values
		}
		m.Instances[inst.Id] = inst
		return inst
	}
	i := makeInstance("i", makeClass("test class", "test classes"))
	makeInstance("x", makeClass("test other", "test others"))
	mem := NewMemoryModel(m, v, make(table.Tables))
	// test plural
	if cls, ok := mem.GetClass("TestClasses"); assert.True(t, ok, "get test class") {
		if p, ok := cls.GetProperty("Plural"); assert.True(t, ok) {
			require.EqualValues(t, i.Class.Plural, p.GetValue().GetText())
		}
		if p, ok := cls.GetProperty("Singular"); assert.True(t, ok) {
			require.EqualValues(t, i.Class.Singular, p.GetValue().GetText())
		}
	}
	// test the api
	tests.ApiTest(t, mem)
	// test that things really changed
	if res, ok := v.GetValue(i.Id, numProp); assert.True(t, ok) {
		require.EqualValues(t, float32(32), res)
	}
	if res, ok := v.GetValue(i.Id, textProp); assert.True(t, ok) {
		require.EqualValues(t, "text", res)
	}
	if res, ok := v.GetValue(i.Id, stateProp); assert.True(t, ok) {
		if c, e := enum.IndexToChoice(res.(int)); assert.NoError(t, e) {
			require.EqualValues(t, ident.Id("Yes"), c)
		}
	}
	if res, ok := v.GetValue(i.Id, objectProp); assert.True(t, ok) {
		require.EqualValues(t, i.Id, res)
	}
}
