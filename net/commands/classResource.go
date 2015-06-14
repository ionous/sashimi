package commands

import (
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/net/resource"
)

func ClassResource(model *M.Model) resource.IResource {
	return resource.Wrapper{
		// find the named class
		Finds: func(id string) (ret resource.IResource, okay bool) {
			if cls, ok := model.Classes.FindClass(id); ok {
				okay, ret = true, resource.Wrapper{
					// return information about the id'd class
					Queries: func(doc resource.DocumentBuilder) {
						addClass(doc, doc.NewIncludes(), cls)
					},
					// find a sub-resource of the id'd class
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
			for k, _ := range model.Classes {
				objects.NewObject(jsonId(k), "class")
			}
		},
	}
}

func classParents(ar []string, cls *M.ClassInfo) []string {
	if p := cls.Parent(); p != nil {
		ar = append(classParents(ar, p), jsonId(p.Id()))
	}
	return ar
}

func addClass(doc, sub resource.IBuildObjects, cls *M.ClassInfo) {
	var parent *resource.Object
	if p := cls.Parent(); p != nil {
		//addClass(sub, sub, p)
		// disabling recursion
		parent = resource.NewObject(jsonId(p.Id()), "class")
	}
	id := jsonId(cls.Id())
	out := doc.NewObject(id, "class")
	out.SetAttr("parent", parent)
	out.SetAttr("name", cls.Name())
	out.SetAttr("singular", cls.Singular())
	out.SetMeta("classes", classParents([]string{id}, cls))
	props := resource.Dict{}
	for pid, prop := range cls.Properties() {
		typeName := "unknown"
		switch prop.(type) {
		case *M.RelativeProperty:
			typeName = "rel"
		case *M.TextProperty:
			typeName = "text"
		case *M.NumProperty:
			typeName = "num"
		case *M.EnumProperty:
			typeName = "enum"
		}
		props[jsonId(pid)] = typeName
	}
	out.SetAttr("props", props)
}
