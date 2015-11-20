package metaltest

import (
	"fmt"
	"github.com/ionous/sashimi/compiler/metal"
	"github.com/ionous/sashimi/compiler/model/modeltest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func PostConditions(t *testing.T, v metal.ObjectValue) {
	// test that things really changed
	if res, ok := v.GetValue(modeltest.TestInstance, modeltest.NumProp); !assert.True(t, ok, fmt.Sprintf("missing (%s.%s)", modeltest.TestInstance, modeltest.NumProp)) {
		t.Fatal()
	} else {
		require.EqualValues(t, float32(32), res)
	}
	if res, ok := v.GetValue(modeltest.TestInstance, modeltest.TextProp); assert.True(t, ok) {
		require.EqualValues(t, "text", res)
	}
	if res, ok := v.GetValue(modeltest.TestInstance, modeltest.StateProp); assert.True(t, ok) {
		require.EqualValues(t, "yes", res)
	}
	if res, ok := v.GetValue(modeltest.TestInstance, modeltest.ObjectProp); assert.True(t, ok) {
		require.EqualValues(t, modeltest.TestInstance, res)
	}
	//
	if res, ok := v.GetValue(modeltest.TestInstance, modeltest.NumsProp); assert.True(t, ok) {
		require.Contains(t, res, float32(32))
	}
	if res, ok := v.GetValue(modeltest.TestInstance, modeltest.TextsProp); assert.True(t, ok) {
		require.Contains(t, res, "text")
	}
	if res, ok := v.GetValue(modeltest.TestInstance, modeltest.ObjectProp); assert.True(t, ok) {
		require.Contains(t, res, modeltest.TestInstance)
	}
}
