package app

import (
	"fmt"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/net/resource"
	"github.com/ionous/sashimi/util/ident"
)

func ObjectResource(mdl meta.Model, clsId ident.Id) resource.IResource {
	return resource.Wrapper{
		// Find the id object.
		Finds: func(name string) (ret resource.IResource, okay bool) {
			id := ident.MakeId(name)
			if inst, ok := mdl.GetInstance(id); ok {
				if cls := inst.GetParentClass(); clsId == cls.GetId() {
					okay, ret = true, resource.Wrapper{
						// Return the object:
						Queries: func(doc resource.DocumentBuilder) {
							doc.AddObject(SerializeObject(inst))
						},
						// Find a relation in the object:
						Finds: func(propertyName string) (ret resource.IResource, okay bool) {
							// FIX: relations are stored in the model
							prop, ok := inst.FindProperty(propertyName)
							if !ok {
								propId := ident.MakeId(propertyName)
								prop, ok = inst.GetProperty(propId)
							}
							if ok {
								if _, ok := prop.GetRelative(); ok {
									okay, ret = true, resource.Wrapper{
										// Return the list of related objects:
										Queries: func(doc resource.DocumentBuilder) {
											objects, includes :=
												NewRefSerializer(mdl, doc.NewObjects()),
												NewObjSerializer(mdl, doc.NewIncludes())

											addObject := func(n ident.Id) {
												if gobj, ok := mdl.GetInstance(n); !ok {
													panic(fmt.Sprintf("internal error, couldnt find related object '%s'", n))
												} else {
													objects.AddObject(gobj)
													includes.SerializeObject(gobj)
												}
											}

											// UGH. for backwards compatibility (ex. whereabouts queries
											if propType := prop.GetType(); propType&meta.ArrayProperty == 0 {
												n := prop.GetValue().GetObject()
												addObject(n)
											} else {
												vals := prop.GetValues()
												for i := 0; i < vals.NumValue(); i++ {
													n := vals.ValueNum(i).GetObject()
													addObject(n)
												}
											}
										},
									}
								}
							}
							return ret, okay
						},
					}
				}
			}
			return ret, okay
		},
	}
}

/* json api for list of objects:
	{ "data": [{
    "type": "articles",
    "id": "1",
    "attributes": {
      "title": "JSON API paints my bikeshed!"
    }
  }, {
    "type": "articles",
    "id": "2",
    "attributes": {
      "title": "Rails is Omakase"
    }
  }]
  }*/
