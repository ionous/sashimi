package tests

import (
	. "github.com/ionous/mars/core"
	. "github.com/ionous/mars/script"
	"github.com/ionous/sashimi/compiler"
	M "github.com/ionous/sashimi/compiler/model"
	"github.com/ionous/sashimi/util/ident"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"sort"
	"testing"
)

// TestRelationRules to see if the relation rules build properly.
func TestRelationRules(t *testing.T) {
	s := buildRelationRules()
	if src, err := s.BuildStatements(); assert.NoError(t, err, "build") {
		if model, err := compiler.Compile(src, ioutil.Discard); assert.NoError(t, err, "compile") {
			if assert.Equal(t, 1, len(model.Relations)) {
				for _, v := range model.Relations {
					require.EqualValues(t, "GremlinsPets", v.Source.String())
					require.EqualValues(t, "RocksOBeneficentOne", v.Target.String())
					require.EqualValues(t, M.OneToMany, v.Style)
				}
				return
			}
		}
	}
	t.FailNow()
}

// TestRelationData to see if the relation values build properly.
func TestRelationData(t *testing.T) {
	s := append(buildRelationRules(), buildRelationValues())
	if src, err := s.BuildStatements(); assert.NoError(t, err, "build") {
		if model, err := compiler.Compile(src, ioutil.Discard); assert.NoError(t, err, "compile") {
			if assert.NoError(t, err, "compile") {
				if assert.Equal(t, 2, len(model.Instances), "two instances") {
					claire, ok := model.Instances[ident.MakeId("claire")]
					require.True(t, ok, "found claire")

					gremlins, ok := model.Classes[claire.Class]
					require.True(t, ok, "found gremlins")

					petsrel, ok := gremlins.FindProperty("pets")
					require.True(t, ok)
					require.True(t, !petsrel.Relation.Empty())

					loofah, ok := model.Instances[ident.MakeId("Loofah")]
					require.True(t, ok, "found loofah")

					rocks, ok := model.Classes[loofah.Class]
					require.True(t, ok, "found rocks")

					gremlinrel, ok := rocks.FindProperty("o beneficent one")
					require.True(t, ok, "found benes")

					omygremlin, ok := loofah.Values[gremlinrel.Id]
					require.True(t, ok, "found grem")
					require.EqualValues(t, claire.Id, omygremlin)
				}
				return
			}
		}
	}
	t.FailNow()
}

// TestRelationIteration to see if the relation values run properly.
func TestRelationIteration(t *testing.T) {
	s := append(buildRelationIteration(), buildRelationRules(), buildRelationValues())
	if test, err := NewTestGameSource(t, &s, "no parser", nil); assert.NoError(t, err) {
		//
		if err := test.Game.RunNamedAction("test", Named{"claire test"}); assert.NoError(t, err) {
			if out, err := test.FlushOutput(); assert.NoError(t, err) {
				sort.Strings(out)
				expected := lines("Boomba", "Loofah")
				sort.Strings(expected)
				require.EqualValues(t, expected, out)
			}
		}
		require.NoError(t, test.Game.RunNamedAction("test", Named{"loofah test"}))
		require.NoError(t, test.Game.RunNamedAction("test", Named{"boomba test"}))
		return
	}
	t.FailNow()
}

func buildRelationIteration() Script {
	return Script{
		The("gremlin", Called("Claire"), Has("pets", "Boomba")),
		The("rock", Called("Boomba"), Exists()),
		The("kinds",
			Called("unit tests"),
			Can("test").And("testing").RequiresNothing()),
		The("unit test", Called("loofah test"),
			When("testing").Always(
				Try("from loofah get claire",
					IsText{
						T("Claire"),
						EqualTo,
						PropertyText{
							PropertyRef{
								Named{"loofah"}, "o beneficent one",
							},
							"name",
						},
					},
				),
			)),
		//
		The("unit test", Called("boomba test"),
			When("testing").Always(
				Try("from boomba get claire",
					IsText{
						T("Claire"),
						EqualTo,
						PropertyText{
							PropertyRef{
								Named{"boomba"}, "o beneficent one",
							},
							"name",
						},
					},
				),
			)),
		// from claire get a list which includes loofah and boomba
		The("unit test", Called("claire test"),
			When("testing").Always(
				EachObj{
					In: PropertyRefList{
						Named{"claire"}, "pets",
					},
					Go: PrintLine{
						PrintText{
							GetText{"name"},
						},
					},
					Else: Error{T("should have run")},
				})),
	}
}

func buildRelationValues() Script {
	return Script{
		The("gremlin", Called("Claire"), Has("pets", "Loofah")),
		The("rock", Called("Loofah"), Exists()),
	}
}

func buildRelationRules() Script {
	return Script{
		The("kinds",
			Called("gremlins"),
			HaveMany("pets", "rocks").
				Implying("rocks", HaveOne("o beneficent one", "gremlin")),
			// alternate, non-conflicting specification of the same relation
			HaveMany("pets", "rocks").
				// FIX: for now the property names must match.
				// FIX? if the names don't match, this creates two views of the same relation. validate the hierarchy to verify no duplicate property usage?
				// i'd prefer the singular: Has("pet", "Loofah")
				Implying("rocks", HaveOne("o beneficent one", "gremlin")),
		),
		The("kinds", Called("rocks"), Exist()),
	}
}
