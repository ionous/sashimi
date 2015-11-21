package tests

import (
	M "github.com/ionous/sashimi/compiler/model"
	. "github.com/ionous/sashimi/script"
	"github.com/ionous/sashimi/util/ident"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

//
// create a single subclass called stories
func TestRelation(t *testing.T) {
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
	res, err := s.Compile(Log(t))
	if assert.NoError(t, err) {
		model := res.Model
		model.PrintModel(t.Log)
		assert.Equal(t, 1, len(model.Relations))
		for _, v := range model.Relations {
			assert.EqualValues(t, "GremlinsPets", v.Source.String())
			assert.EqualValues(t, "RocksOBeneficentOne", v.Target.String())
			assert.EqualValues(t, M.OneToMany, v.Style)
		}
	}
}

//
// create a single subclass called stories
func TestRelates(t *testing.T) {
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
	res, err := s.Compile(Log(t))
	if assert.NoError(t, err, "compile") {
		model := res.Model
		assert.Equal(t, 2, len(model.Instances), "two instances")

		claire, ok := model.Instances[ident.MakeId("claire")]
		require.True(t, ok, "found claire")

		gremlins, ok := model.Classes[claire.Class]
		require.True(t, ok, "found gremlins")

		petsrel, ok := gremlins.FindProperty("pets")
		assert.True(t, ok)
		assert.True(t, !petsrel.Relation.Empty())

		loofah, ok := model.Instances[ident.MakeId("Loofah")]
		assert.True(t, ok, "found loofah")

		rocks, ok := model.Classes[loofah.Class]
		require.True(t, ok, "found rocks")

		gremlinrel, ok := rocks.FindProperty("o beneficent one")
		require.True(t, ok, "found benes")

		omygremlin, ok := loofah.Values[gremlinrel.Id]
		require.True(t, ok, "found grem")
		assert.EqualValues(t, claire.Id, omygremlin)
	}
}
