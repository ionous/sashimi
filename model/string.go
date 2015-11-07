package model

import (
	"github.com/ionous/sashimi/util/ident"
	"github.com/ionous/sashimi/util/lang"
)

// FIX? pretty much everything that has a map of ident.Id also needs a find by name
type StringPair struct {
	id  ident.Id
	str string
}

func MakeStringId(name string) ident.Id {
	return ident.MakeId(lang.StripArticle(name))
}
