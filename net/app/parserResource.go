package app

import (
	"github.com/ionous/sashimi/net/resource"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/util/ident"
)

func ParserResource(m api.Model) resource.IResource {
	return resource.Wrapper{
		Queries: func(doc resource.DocumentBuilder) {
			objects := doc.NewObjects()
			for i := 0; i < m.NumParserAction(); i++ {
				act := m.ParserActionNum(i)
				objects.NewObject(jsonId(act.GetId()), "action")
			}
		},
		Finds: func(id string) (ret resource.IResource, okay bool) {
			if act, ok := m.GetAction(ident.MakeId(id)); ok {
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
