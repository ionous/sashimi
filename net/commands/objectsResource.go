package commands

import (
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/net/resource"
	R "github.com/ionous/sashimi/runtime"
)

func ClassRuntime(game *R.Game, cls *M.ClassInfo, serial *ObjectSerializer) resource.IResource {
	return resource.Wrapper{
		// Find the id object.
		Finds: func(id string) (ret resource.IResource, okay bool) {
			if gobj, ok := game.Objects[M.MakeStringId(id)]; ok && gobj.Class() == cls {
				okay, ret = true, resource.Wrapper{
					// Return the object:
					Queries: func(doc resource.DocumentBuilder) {
						serial.SerializeObject(doc, gobj, true)
					},
					// Find a relation in the object:
					Finds: func(propertyName string) (ret resource.IResource, okay bool) {
						// FIX: relations are stored in the model
						inst := game.Model.Instances[gobj.Id()]
						if rel, ok := inst.RelativeValue(propertyName); ok {
							okay, ret = true, resource.Wrapper{
								// Return the list of related objects:
								Queries: func(doc resource.DocumentBuilder) {
									classes, includes := doc.NewObjects(), doc.NewIncludes()
									//
									for _, n := range rel.List() {
										gobj := game.Objects[M.StringId(n)]
										serial.AddObjectRef(classes, gobj, includes)
									}
								},
							}
						}
						return ret, okay
					},
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
  }
*/
