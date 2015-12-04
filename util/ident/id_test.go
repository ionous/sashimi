package ident

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBasics(t *testing.T) {
	titleCase := MakeId("TitleCase")
	assert.Equal(t, "TitleCase", titleCase.String(), "title case stays tile case")
	assert.Equal(t, "TwoWords", MakeId("two words").String(), "TwoWords", "two words become one")
	assert.Equal(t, "WordDash", MakeId("word-dash").String(), "WordDash", "dashed words split")
	assert.Equal(t, "Apostrophes", MakeId("apostrophe's").String(), "apostrophes vanish")
	assert.True(t, Empty().Empty(), "empty is as empty does")
	assert.False(t, MakeUniqueId().Empty(), "unique works")
	assert.Equal(t, "TitleCase", Join(MakeId("title"), MakeId("case")).String(), "title case stays tile case")
	assert.Equal(t, MakeId("786abc123def"), MakeId("786-abc 123 def"))
}

func TestCompare(t *testing.T) {
	assert.Equal(t, 0, Compare(MakeId("a b c"), MakeId("A B C")))
	assert.Equal(t, -1, Compare(MakeId("a b"), MakeId("A B C")))
	assert.Equal(t, 1, Compare(MakeId("a b c"), MakeId("A B")))
	assert.Equal(t, 1, Compare(MakeId("a b c d"), MakeId("A B C")))
	assert.Equal(t, 1, Compare(MakeId("aba"), MakeId("ab")))
	assert.Equal(t, -1, Compare(MakeId("ab"), MakeId("aba")))
	assert.Equal(t, 0, Compare(MakeId("aba"), MakeId("aba")))
	//return-policy-1 ReturnPolicy1 return-policy1 ReturnPolicy1 false
	// id2 := ident.MakeId("ReturnPolicy1")
	assert.Equal(t, 0, Compare(MakeId("return policy 1"), MakeId("ReturnPolicy1")))
}

func TestToString(t *testing.T) {
	assert.EqualValues(t, "Apples", MakeId("apples").String())
	assert.EqualValues(t, "Apples", MakeId("Apples").String())

	assert.EqualValues(t, "AppleTurnover", MakeId("apple turnover").String())
	assert.EqualValues(t, "AppleTurnover", MakeId("Apple Turnover").String())
	assert.EqualValues(t, "AppleTurnover", MakeId("Apple turnover").String())
	assert.EqualValues(t, "AppleTurnover", MakeId("APPLE TURNOVER").String())

	assert.EqualValues(t, "PascalCase", MakeId("PascalCase").String())
	assert.EqualValues(t, "CamelCase", MakeId("camelCase").String())

	assert.EqualValues(t, "SomethingLikeThis", MakeId("something-like-this").String())
	assert.EqualValues(t, "Allcaps", MakeId("ALLCAPS").String())

	assert.EqualValues(t, "WhaTaboutThis", MakeId("whaTAboutThis").String())

	assert.Exactly(t, []string{"apple", "turnover"}, MakeId("AppleTurnover").Split())
	assert.Exactly(t, []string{"title"}, MakeId("Title").Split())
}

func TestDash(t *testing.T) {
	assert.EqualValues(t, "apple-turnover", Dash(MakeId("apple turnover")))
	assert.EqualValues(t, "dash-identity", Dash(MakeId("dash-identity")))
	assert.EqualValues(t, "dash-it-all", Dash(MakeId("DashItAll")))
}
