package resource

import "fmt"

type ObjectsBuilder struct {
	builder DocumentBuilder
}

//
// Add an object identifier to the list of document objects.
// Return a builder to turn the identifier into a full object.
//
func (this ObjectsBuilder) NewObject(id, class string) *Object {
	obj := NewObject(id, class)
	switch data := this.builder.doc.Data.(type) {
	// our first object was a blank placeholder:
	case []Object:
		this.builder.doc.Data = []*Object{obj}
	// otherwise we have a pointer array of modifying the objects
	case []*Object:
		this.builder.doc.Data = append(data, obj)
	default:
		this.builder.AddError(fmt.Errorf("multiple objects added to a single object document."))
	}
	return obj
}
