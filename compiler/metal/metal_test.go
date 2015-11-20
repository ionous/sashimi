package metal

import (
	"fmt"
	M "github.com/ionous/sashimi/compiler/model/modeltest"
	"github.com/ionous/sashimi/compiler/model/table"
	"github.com/ionous/sashimi/meta/tests"
	"github.com/ionous/sashimi/util/ident"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// go test -run TestMetal
func TestMetal(t *testing.T) {
	v := make(ObjectValueMap)
	m := NewMetal(M.NewModel(), v, make(table.Tables))

	// test plural
	if cls, ok := m.GetClass(ident.MakeId("TestClasses")); assert.True(t, ok, "get test class") {
		if p, ok := cls.FindProperty("plural"); assert.True(t, ok) {
			require.EqualValues(t, "test classes", p.GetValue().GetText())
		}
		if p, ok := cls.FindProperty("singular"); assert.True(t, ok) {
			require.EqualValues(t, "test class", p.GetValue().GetText())
		}
	}
	// test the api
	tests.ApiTest(t, m, M.TestInstance)
	// test that things really changed
	if res, ok := v.GetValue(M.TestInstance, M.NumProp); !assert.True(t, ok, fmt.Sprintf("missing (%s.%s)", M.TestInstance, M.NumProp)) {
		t.Log(v)
		t.Fatal()
	} else {
		require.EqualValues(t, float32(32), res)
	}
	if res, ok := v.GetValue(M.TestInstance, M.TextProp); assert.True(t, ok) {
		require.EqualValues(t, "text", res)
	}
	if res, ok := v.GetValue(M.TestInstance, M.StateProp); assert.True(t, ok) {
		// 1= No, 2= Yes.
		// Why *does** GetValue return an integer anyway?
		require.EqualValues(t, 2, res)

	}
	if res, ok := v.GetValue(M.TestInstance, M.ObjectProp); assert.True(t, ok) {
		require.EqualValues(t, M.TestInstance, res)
	}
	//
	if res, ok := v.GetValue(M.TestInstance, M.NumsProp); assert.True(t, ok) {
		require.Contains(t, res, float32(32))
	}
	if res, ok := v.GetValue(M.TestInstance, M.TextsProp); assert.True(t, ok) {
		require.Contains(t, res, "text")
	}
	if res, ok := v.GetValue(M.TestInstance, M.ObjectProp); assert.True(t, ok) {
		require.Contains(t, res, M.TestInstance)
	}
}
