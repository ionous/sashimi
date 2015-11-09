package app

import (
	"fmt"
	"github.com/ionous/sashimi/net/resource"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/util/ident"
)

func ObjectResource(mdl api.Model, clsId ident.Id, serial *ObjectSerializer) resource.IResource {
	return resource.Wrapper{
		// Find the id object.
		Finds: func(name string) (ret resource.IResource, okay bool) {
			id := ident.MakeId(name)
			if inst, ok := mdl.GetInstance(id); ok {
				if cls := inst.GetParentClass(); clsId == cls.GetId() {
					okay, ret = true, resource.Wrapper{
						// Return the object:
						Queries: func(doc resource.DocumentBuilder) {
							serial.SerializeObject(doc, inst, true)
						},
						// Find a relation in the object:
						Finds: func(propertyName string) (ret resource.IResource, okay bool) {
							// FIX: relations are stored in the model
							propId := ident.MakeId(propertyName)

							if prop, ok := inst.GetProperty(propId); ok {
								if prop.GetType() == api.ObjectProperty|api.ArrayProperty {
									vals := prop.GetValues()
									okay, ret = true, resource.Wrapper{
										// Return the list of related objects:
										Queries: func(doc resource.DocumentBuilder) {
											classes, includes := doc.NewObjects(), doc.NewIncludes()
											//
											for i := 0; i < vals.NumValue(); i++ {
												n := vals.ValueNum(i).GetObject()
												if other, ok := mdl.GetInstance(n); !ok {
													panic(fmt.Sprintf("internal error, couldnt find related object '%s'", n))
												} else {
													serial.AddObjectRef(classes, other, includes)
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
