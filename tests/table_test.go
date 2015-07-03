package tests

import (
	M "github.com/ionous/sashimi/model"
	. "github.com/ionous/sashimi/script"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func sweets(s *Script) *Script {
	s.The("kinds",
		Called("sweets"),
		Have("desc", "text"),
		AreOneOf("delicious", "decent", "acceptable", "you can't be serious"))
	s.The("sweet",
		Table("name", "desc", "delicious-property").Contains(
			"boreo", "creme filled wafer things.", "acceptable").And(
			"uncle eddie's vegan chocolate chip cookies", "the name says it all.", "delicious").And(
			"sugar coated ants", "ants. coated with sugar.", "you can't be serious",
		),
	)
	return s
}

// 1. create a simple table declaration and generate some fixed instances
// 2. create a simple table declaration with autogenerated instances
// 3. merge some non-contrary instance data
// 4. test table declarations which use pointers
func TestTableDecl(t *testing.T) {
	s := sweets(&Script{})
	if m, err := s.Compile(os.Stderr); assert.NoError(t, err, "table compile") {
		assert.Len(t, m.Instances, 3)
		// and now some values:
		if inst, ok := m.Instances.FindInstance("boreo"); assert.True(t, ok, "find tabled instance by name") {
			//
			if val, ok := inst.ValueByName("desc"); assert.True(t, ok, "find desc") {
				if str, ok := val.(*M.TextValue); assert.True(t, ok, "have string") {
					assert.EqualValues(t, str.String(), "creme filled wafer things.")
				}
			}
			//
			if val, ok := inst.ValueByName("delicious-property"); assert.True(t, ok, "find deliciousness") {
				if enum, ok := val.(*M.EnumValue); assert.True(t, ok, "have enum") {
					assert.EqualValues(t, enum.String(), "acceptable")
				}
			}
			//
			if val, ok := inst.ValueByName("delicious-property"); assert.True(t, ok, "find deliciousness") {
				if enum, ok := val.(*M.EnumValue); assert.True(t, ok, "have enum") {
					assert.EqualValues(t, enum.String(), "acceptable")
				}
			}
		}
	}
}
