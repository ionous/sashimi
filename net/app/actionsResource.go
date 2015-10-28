package app

import (
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/net/resource"
)

// func ActionsResource(model *M.Model) resource.IResource {
// 	return resource.Wrapper{
// 		Finds: func(name string) (ret resource.IResource, okay bool) {
// 			if act, e := model.Actions.FindActionByName(name); e == nil {
// 				okay, ret = true, resource.Wrapper{
// 					// action data:
// 					Queries: func(doc resource.DocumentBuilder) {
// 						actionResource(doc, act)
// 					},
// 				}
// 			}
// 			return ret, okay
// 		}}
// }

func actionResource(out resource.IBuildObjects, act *M.ActionInfo) {
	out.NewObject(jsonId(act.Id), "action").
		SetAttr("act", act.ActionName).
		SetAttr("evt", act.EventName).
		SetAttr("src", classToId(act.Source())).
		SetAttr("tgt", classToId(act.Target())).
		SetAttr("ctx", classToId(act.Context()))
}

func classToId(cls *M.ClassInfo) (ret string) {
	if cls != nil {
		ret = jsonId(cls.Id)
	}
	return ret
}
