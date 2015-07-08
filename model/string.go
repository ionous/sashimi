package model

import (
	"github.com/ionous/sashimi/util/ident"
	"regexp"
)

//
// FIX: pretty much everything that has a map of ident.Id also needs a find by name
//
type StringPair struct {
	id  ident.Id
	str string
}

//
// unfortunately, the parser has this too.
// FIX: move to some common space?
//
var articles = regexp.MustCompile(`^(The|A|An|Our|Some)[[:upper:]]`)

func stripArticle(name string) (article, bare string) {
	pair := articles.FindStringIndex(name)
	if pair == nil {
		bare = name
	} else {
		split := pair[1] - 1
		article = name[:split]
		bare = name[split:]
	}
	return article, bare
}

func MakeStringId(name string) ident.Id {
	_, str := stripArticle(ident.MakeId(name).String())
	return ident.Id(str)
}
