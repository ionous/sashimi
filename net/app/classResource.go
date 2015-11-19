package app

import (
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/net/resource"
	"github.com/ionous/sashimi/util/ident"
)

func ClassResource(mdl meta.Model) resource.IResource {
	return resource.Wrapper{
		// list all classes
		Queries: func(doc resource.DocumentBuilder) {
			objects := doc.NewObjects()
			for i := 0; i < mdl.NumClass(); i++ {
				cls := mdl.ClassNum(i)
				objects.NewObject(jsonId(cls.GetId()), "class")
			}
		},
		// find the named class
		Finds: func(id string) (ret resource.IResource, okay bool) {
			if cls, ok := mdl.GetClass(ident.MakeId(id)); ok {
				okay, ret = true, resource.Wrapper{
					// return information about the id'd class
					Queries: func(doc resource.DocumentBuilder) {
						addClass(doc, doc.NewIncludes(), cls)
					},
				}
			}
			return
		},
	}
}

func classParents(cls meta.Class, ar []string) []string {
	if p := cls.GetParentClass(); p != nil {
		ar = append(classParents(p, ar), jsonId(p.GetId()))
	}
	return ar
}

func addClass(doc, sub resource.IBuildObjects, cls meta.Class) {
	var parent *resource.Object
	if p := cls.GetParentClass(); p != nil {
		//addClass(model, sub, sub, p)
		// disabling recursion
		parent = resource.NewObject(jsonId(p.GetId()), "class")
	}
	id := jsonId(cls.GetId())
	out := doc.NewObject(id, "class")
	out.SetAttr("parent", parent)
	plural, _ := cls.FindProperty("plural")
	out.SetAttr("name", plural.GetValue().GetText())
	singular, _ := cls.FindProperty("singular")
	out.SetAttr("singular", singular.GetValue().GetText())
	a := append(classParents(cls, nil), id)
	// reverse
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
	out.SetMeta("classes", a)
	props := resource.Dict{}
	for i := 0; i < cls.NumProperty(); i++ {
		prop := cls.PropertyNum(i)
		typeName := "unknown"
		switch prop.GetType() {
		case meta.ArrayProperty | meta.ObjectProperty:
			typeName = "rel"
		case meta.TextProperty:
			typeName = "text"
		case meta.NumProperty:
			typeName = "num"
		case meta.StateProperty:
			typeName = "enum"
		}
		props[jsonId(prop.GetId())] = typeName
	}
	// actionRefs := resource.NewObjectList()
	// for _, act := range model.Actions {
	// 	var dst ident.Id
	// 	if d := act.Context(); d != nil {
	// 		dst = d.Id()
	// 	} else if d := act.Target(); d != nil {
	// 		dst = d.Id()
	// 	}
	// 	if cls.CompatibleWith(dst) {
	// 		actionRefs.NewObject(jsonId(act.Id()), "action")
	// 		actionResource(sub, act)
	// 	}
	// }
	out.SetAttr("props", props)
	// /out.SetAttr("actions", actionRefs.Objects())
}
