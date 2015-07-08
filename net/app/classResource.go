package app

import (
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/net/resource"
	"github.com/ionous/sashimi/util/ident"
)

func ClassResource(model *M.Model) resource.IResource {
	return resource.Wrapper{
		// list all classes
		Queries: func(doc resource.DocumentBuilder) {
			objects := doc.NewObjects()
			for k, _ := range model.Classes {
				objects.NewObject(jsonId(k), "class")
			}
		},
		// find the named class
		Finds: func(id string) (ret resource.IResource, okay bool) {
			if cls, ok := model.Classes.FindClass(id); ok {
				okay, ret = true, resource.Wrapper{
					// return information about the id'd class
					Queries: func(doc resource.DocumentBuilder) {
						addClass(model, doc, doc.NewIncludes(), cls)
					},
				}
			}
			return
		},
	}
}

func classParents(cls *M.ClassInfo, ar []string) []string {
	if p := cls.Parent(); p != nil {
		ar = append(classParents(p, ar), jsonId(p.Id()))
	}
	return ar
}

func addClass(model *M.Model, doc, sub resource.IBuildObjects, cls *M.ClassInfo) {
	var parent *resource.Object
	if p := cls.Parent(); p != nil {
		//addClass(model, sub, sub, p)
		// disabling recursion
		parent = resource.NewObject(jsonId(p.Id()), "class")
	}
	id := jsonId(cls.Id())
	out := doc.NewObject(id, "class")
	out.SetAttr("parent", parent)
	out.SetAttr("name", cls.Name())
	out.SetAttr("singular", cls.Singular())
	a := append(classParents(cls, nil), id)
	// reverse
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
	out.SetMeta("classes", a)
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
	actionRefs := resource.NewObjectList()
	for _, act := range model.Actions {
		var dst ident.Id
		if d := act.Context(); d != nil {
			dst = d.Id()
		} else if d := act.Target(); d != nil {
			dst = d.Id()
		}
		if cls.CompatibleWith(dst) {
			actionRefs.NewObject(jsonId(act.Id()), "action")
			actionResource(sub, act)
		}
	}
	out.SetAttr("props", props)
	out.SetAttr("actions", actionRefs.Objects())
}
