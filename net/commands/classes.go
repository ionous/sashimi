package commands

import (
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/net/resource"
)

func ClassModel(model *M.Model, cls *M.ClassInfo) resource.IResource {
	return resource.Wrapper{
		// find the named class
		Finds: func(id string) (ret resource.IResource, okay bool) {
			if cls, ok := model.Classes.FindClass(id); ok {
				okay, ret = true, resource.Wrapper{
					// find the class sub-resource
					Finds: func(name string) (ret resource.IResource, okay bool) {
						switch name {
						case "actions":
							okay, ret = true, resource.Wrapper{
								Queries: func(doc resource.DocumentBuilder) {
									objects := doc.NewObjects()
									for _, act := range model.Actions {
										if act.Target() == cls || cls.HasParent(act.Target()) {
											objects.NewObject(jsonId(act.Id()), "action")
										}
									}
								}}
						}
						return
					}}
			}
			return
		},
		// list all classes
		Queries: func(doc resource.DocumentBuilder) {
			objects := doc.NewObjects()
			for k, c := range model.Classes {
				objects.NewObject(jsonId(k), "class").
					SetAttr("name", c.Name())
			}
		},
	}
}

// 						//	case "relatives":
// 						// 		okay, ret = true, Wrapper{
// 						// 			Queries: func() interface{} {
// 						// 				relatives := []*M.RelativeProperty{}
// 						// 				for _, v := range cls.AllProperties() {
// 						// 					if rel, ok := v.(*M.RelativeProperty); ok {
// 						// 						relatives = append(relatives, rel)
// 						// 					}
// 						// 				}
// 						// 				return relatives
// 						// 			},
// 						// 		}
// 						// 	}
// 						// 	return
// 						return
// 				}
// 			}
// 			return
// 		},
// 	}
// }
