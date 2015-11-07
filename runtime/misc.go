package runtime

import (
	"github.com/ionous/sashimi/util/ident"
	"github.com/ionous/sashimi/util/lang"
)

func StripStringId(name string) ident.Id {
	return ident.MakeId(lang.StripArticle(name))
}

func MakeStringId(name string) ident.Id {
	return ident.MakeId(name)
}
