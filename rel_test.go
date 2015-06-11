package sashimi

import (
	M "github.com/ionous/sashimi/model"
	. "github.com/ionous/sashimi/script"
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
		Have("pets", "rocks").
			Implying("rocks", Have("o beneficent one", "gremlin")),
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
		Have("pets", "rocks").
			Implying("rocks", Have("o beneficent one", "gremlin")),
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

		claire, err := model.Instances.FindInstance("claire")
		assert.NoError(t, err, "found claire")

		pets, ok := claire.ValueByName("pets")
		assert.True(t, ok)
		petsrel := pets.(*M.RelativeValue)
		assert.Exactly(t, []string{"Loofah"}, petsrel.List())

		loofah, err := model.Instances.FindInstance("loofah")
		assert.NoError(t, err, "found loofah")

		gremlins, ok := loofah.ValueByName("o beneficent one")
		assert.True(t, ok, "value by name")
		gremlinrel := gremlins.(*M.RelativeValue)
		assert.Exactly(t, []string{"Claire"}, gremlinrel.List())

		model.PrintModel(t.Log)
	}
}
