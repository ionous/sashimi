package commands

import (
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/net/resource"
)

func ActionsResource(model *M.Model) resource.IResource {
	return resource.Wrapper{
		Finds: func(name string) (ret resource.IResource, okay bool) {
			if act, e := model.Actions.FindActionByName(name); e == nil {
				okay, ret = true, resource.Wrapper{
					// action data:
					Queries: func(doc resource.DocumentBuilder) {
						doc.NewObject(jsonId(act.Id()), "action").
							SetAttr("act", act.Action()).
							SetAttr("evt", act.Event()).
							SetAttr("src", classToId(act.Source())).
							SetAttr("tgt", classToId(act.Target())).
							SetAttr("ctx", classToId(act.Context()))
					},
				}
			}
			return ret, okay
		}}
}

func classToId(cls *M.ClassInfo) (ret string) {
	if cls != nil {
		ret = jsonId(cls.Id())
	}
	return ret
}
