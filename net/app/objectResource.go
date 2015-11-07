package app

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/net/resource"
	R "github.com/ionous/sashimi/runtime"
	"github.com/ionous/sashimi/runtime/api"
)

func ObjectResource(game *R.Game, gcls *M.ClassInfo, serial *ObjectSerializer) resource.IResource {
	return resource.Wrapper{
		// Find the id object.
		Finds: func(id string) (ret resource.IResource, okay bool) {
			if gobj, ok := game.Objects[M.MakeStringId(id)]; ok {
				if cls := gobj.Class(); cls.GetId() == gcls.Id {
					okay, ret = true, resource.Wrapper{
						// Return the object:
						Queries: func(doc resource.DocumentBuilder) {
							serial.SerializeObject(doc, gobj, true)
						},
						// Find a relation in the object:
						Finds: func(propertyName string) (ret resource.IResource, okay bool) {
							// FIX: relations are stored in the model
							propId := M.MakeStringId(propertyName)
							i, _ := game.ModelApi.GetInstance(gobj.Id())

							if prop, ok := i.GetProperty(propId); ok {
								if prop.GetType() == api.ObjectProperty|api.ArrayProperty {
									vals := prop.GetValues()
									okay, ret = true, resource.Wrapper{
										// Return the list of related objects:
										Queries: func(doc resource.DocumentBuilder) {
											classes, includes := doc.NewObjects(), doc.NewIncludes()
											//
											for i := 0; i < vals.NumValue(); i++ {
												n := vals.ValueNum(i).GetObject()
												if gobj, ok := game.Objects[n]; !ok {
													panic(fmt.Sprintf("internal error, couldnt find related object '%s'", n))
												} else {
													serial.AddObjectRef(classes, gobj, includes)
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
