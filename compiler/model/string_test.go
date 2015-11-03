package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringIds(t *testing.T) {
	assert.EqualValues(t, "Apples", MakeStringId("apples"))
	assert.EqualValues(t, "Apples", MakeStringId("Apples"))

	assert.EqualValues(t, "AppleTurnover", MakeStringId("apple turnover"))
	assert.EqualValues(t, "AppleTurnover", MakeStringId("Apple Turnover"))
	assert.EqualValues(t, "AppleTurnover", MakeStringId("Apple turnover"))
	assert.EqualValues(t, "AppleTurnover", MakeStringId("APPLE TURNOVER"))

	assert.EqualValues(t, "PascalCase", MakeStringId("PascalCase"))
	assert.EqualValues(t, "CamelCase", MakeStringId("camelCase"))
	assert.EqualValues(t, "Article", MakeStringId("the article"))

	assert.EqualValues(t, "SomethingLikeThis", MakeStringId("something-like-this"))
	assert.EqualValues(t, "Allcaps", MakeStringId("ALLCAPS"))

	assert.EqualValues(t, "WhaTaboutThis", MakeStringId("whaTAboutThis"))

	assert.Exactly(t, []string{"apple", "turnover"}, MakeStringId("AppleTurnover").Split())
	assert.Exactly(t, []string{"title"}, MakeStringId("Title").Split())
}
