package tests

import (
	M "github.com/ionous/sashimi/model"
	. "github.com/ionous/sashimi/script"

	"github.com/ionous/sashimi/util/ident"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

//
// create a single subclass called stories
func TestSimpleRelation(t *testing.T) {
	s := Script{}
	s.The("kinds",
		Called("gremlins"),
		HaveMany("pets", "rocks").
			Implying("rocks", HaveOne("o beneficent one", "gremlin")),
		// alternate, non-conflicting specification of the same relation
		HaveMany("pets", "rocks").
			// FIX? if the names don't match, this creates two views of the same relation.
			// validate the hierarchy to verify no duplicate property usage?
			Implying("rocks", HaveOne("o beneficent one", "gremlin")),
	)
	s.The("kinds", Called("rocks"), Exist())
	model, err := s.Compile(os.Stderr)
	if assert.NoError(t, err) {
		model.PrintModel(t.Log)
		assert.Equal(t, 1, len(model.Relations))
		for _, v := range model.Relations {
			assert.EqualValues(t, "Gremlins", v.Source().Class)
			assert.EqualValues(t, "Rocks", v.Destination().Class)
			assert.EqualValues(t, M.OneToMany, v.Style())
		}
	}
}

//
// create a single subclass called stories
func TestSimpleRelates(t *testing.T) {
	s := Script{}
	s.The("kinds",
		Called("gremlins"),
		HaveMany("pets", "rocks").
			Implying("rocks", HaveOne("o beneficent one", "gremlin")),
	)
	s.The("kinds", Called("rocks"), Exist())
	// FIX: for now the property names must match,
	// i'd prefer the signular: Has("pet", "Loofah")
	s.The("gremlin", Called("Claire"), Has("pets", "Loofah"))
	s.The("rock", Called("Loofah"), Exists())

	//
	model, err := s.Compile(os.Stderr)
	if assert.NoError(t, err, "compile") {
		assert.Equal(t, 2, len(model.Instances), "two instances")

		claire, ok := model.Instances.FindInstance("claire")
		assert.True(t, ok, "found claire")

		pets, ok := claire.Class().FindProperty("pets")
		assert.True(t, ok)

		petsrel := pets.(*M.RelativeProperty)
		table, ok := model.Tables.TableById(petsrel.Relation())
		assert.True(t, ok, "found table")

		assert.EqualValues(t, []ident.Id{"Loofah"}, table.List(claire.Id(), petsrel.IsRev()))

		loofah, ok := model.Instances.FindInstance("loofah")
		assert.True(t, ok, "found loofah")

		revField, ok := loofah.Class().FindProperty("o beneficent one")
		assert.True(t, ok, "rev property")
		revProp := revField.(*M.RelativeProperty)

		assert.Exactly(t, []ident.Id{"Claire"}, table.List(loofah.Id(), revProp.IsRev()))

		model.PrintModel(t.Log)
	}
}
