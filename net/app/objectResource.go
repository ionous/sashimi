package app

import (
	"fmt"
	"github.com/ionous/sashimi/net/resource"
	R "github.com/ionous/sashimi/runtime"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/util/ident"
)

func ObjectResource(game *R.Game, clsId ident.Id, serial *ObjectSerializer) resource.IResource {
	return resource.Wrapper{
		// Find the id object.
		Finds: func(id string) (ret resource.IResource, okay bool) {
			if gobj, ok := game.Objects[R.MakeStringId(id)]; ok {
				if cls := gobj.Class(); clsId == cls.GetId() {
					okay, ret = true, resource.Wrapper{
						// Return the object:
						Queries: func(doc resource.DocumentBuilder) {
							serial.SerializeObject(doc, gobj, true)
						},
						// Find a relation in the object:
						Finds: func(propertyName string) (ret resource.IResource, okay bool) {
							// FIX: relations are stored in the model
							propId := R.MakeStringId(propertyName)
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
