package app

import (
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/net/resource"
	"github.com/ionous/sashimi/util/ident"
)

func ParserResource(mdl meta.Model) resource.IResource {
	return resource.Wrapper{
		Queries: func(doc resource.DocumentBuilder) {
			objects := doc.NewObjects()
			for i := 0; i < mdl.NumParserAction(); i++ {
				act := mdl.ParserActionNum(i)
				objects.NewObject(jsonId(act.Action), "action")
			}
		},
		Finds: func(id string) (ret resource.IResource, okay bool) {
			if act, ok := mdl.GetAction(ident.MakeId(id)); ok {
				okay, ret = true, resource.Wrapper{
					Queries: func(doc resource.DocumentBuilder) {
						actionResource(doc, act)
					},
				}
			}
			return
		},
	}
}
