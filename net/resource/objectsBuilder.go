package resource

import "fmt"

type ObjectsBuilder struct {
	builder DocumentBuilder
}

// Add an object identifier to the list of document objects.
// Return a builder to turn the identifier into a full object.
func (o ObjectsBuilder) NewObject(id, class string) *Object {
	obj := NewObject(id, class)
	o.AddObject(obj)
	return obj
}

func (o ObjectsBuilder) AddObject(obj *Object) {
	switch data := o.builder.doc.Data.(type) {
	// our first object was a blank placeholder:
	case []Object:
		o.builder.doc.Data = []*Object{obj}
	// otherwise we have a pointer array of modifying the objects
	case []*Object:
		o.builder.doc.Data = append(data, obj)
	default:
		o.builder.AddError(fmt.Errorf("multiple objects added to a single object document."))
	}
}
