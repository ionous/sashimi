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
			m.GetParserActions(func(act api.Action, _ []string) (finished bool) {
				objects.NewObject(jsonId(act.GetId()), "action")
				return finished
			})
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
