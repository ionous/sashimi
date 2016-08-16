package metaltest

import (
	"github.com/ionous/sashimi/compiler/model/modeltest"
	"github.com/ionous/sashimi/meta/metatest"
	"github.com/ionous/sashimi/metal"
	"github.com/ionous/sashimi/util/ident"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// go test -run TestMetal
func TestMetal(t *testing.T) {
	src := make(metal.ObjectValueMap)
	m := metal.NewMetal(modeltest.NewModel(), src)

	// test plural
	if cls, ok := m.GetClass(ident.MakeId("TestClasses")); assert.True(t, ok, "get test class") {
		if p, ok := cls.FindProperty("plural"); assert.True(t, ok) {
			require.EqualValues(t, "test classes", p.GetValue().GetText(), "plural test")
		}
		if p, ok := cls.FindProperty("singular"); assert.True(t, ok) {
			require.EqualValues(t, "test class", p.GetValue().GetText(), "singular test")
		}
	}

	metatest.ApiTest(t, m, modeltest.TestInstance)
	VerifyPostConditions(t, src)
}
