package model_test

import (
	"encoding/json"
	M "github.com/ionous/sashimi/compiler/model"
	"github.com/ionous/sashimi/compiler/model/modeltest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// go test -run TestModelSerial
func TestModelSerial(t *testing.T) {
	src := modeltest.NewModel()
	if b, e := json.Marshal(src); assert.NoError(t, e, "marshal") {
		var retcon *M.Model
		if e := json.Unmarshal(b, &retcon); assert.NoError(t, e, "remarshal") {
			// why does len reverse expected and value?
			assert.Len(t, retcon.Actions, len(src.Actions), "actions")
			assert.Len(t, retcon.Classes, len(src.Classes), "classes")
			assert.Len(t, retcon.Enumerations, len(src.Enumerations), "enums")
			assert.Len(t, retcon.Events, len(src.Events), "events")
			assert.Len(t, retcon.Instances, len(src.Instances), "instances")
			assert.Len(t, retcon.Aliases, len(src.Aliases), "aliases")
			assert.Len(t, retcon.ParserActions, len(src.Actions), "actions")
			assert.Len(t, retcon.Relations, len(src.Relations), "relations")
			assert.Len(t, retcon.SingleToPlural, len(src.SingleToPlural), "plurals")
			assert.EqualValues(t, src, retcon, "deep equal")
			if x, e := json.Marshal(retcon); assert.NoError(t, e, "remarshal") {
				require.EqualValues(t, b, x, "first serialization should match second serialization")
			}
		}
	}
}
