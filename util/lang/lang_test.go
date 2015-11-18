package lang

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStripArticle(t *testing.T) {
	type Pair struct {
		src, article, text string
	}
	p := []Pair{
		{src: "the evil fish", article: "the", text: "evil fish"},
		{src: "The Capital", article: "The", text: "Capital"},
		{src: "some fish", article: "some", text: "fish"},
		{src: " a   space ", article: "a", text: "space"},
		{src: "dune, a desert planet", article: "", text: "dune, a desert planet"},
	}

	for _, p := range p {
		a, b := SliceArticle(p.src)
		require.Equal(t, p.text, b, fmt.Sprintf("text: '%s'", p.src))
		require.Equal(t, p.article, a, fmt.Sprintf("text: '%s'", p.src))
	}

	require.EqualValues(t, "article", StripArticle("the article"))
}
