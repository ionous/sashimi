package app

import (
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/net/resource"
)

func ParserResource(m *M.Model) resource.IResource {
	return resource.Wrapper{
		Queries: func(doc resource.DocumentBuilder) {
			objects := doc.NewObjects()
			for _, act := range m.ParserActions {
				k := act.Action.Id
				objects.NewObject(jsonId(k), "action")
			}
		},
		Finds: func(id string) (ret resource.IResource, okay bool) {
			if act, ok := m.Actions.FindActionByName(id); ok {
				okay, ret = true, resource.Wrapper{
					// return information about the id'd class
					Queries: func(doc resource.DocumentBuilder) {
						actionResource(doc, act)
					},
				}
			}
			return
		},
	}
}
