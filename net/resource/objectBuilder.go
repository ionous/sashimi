package resource

// DocumentBuilder, ObjectsBuilder, and ObjectList satisfy this interface,
// which allows for either one object or multiple objects to be added depending on the context.
type IBuildObjects interface {
	NewObject(id, class string) *Object
	AddObject(obj *Object)
}

type ObjectBuilder struct {
	Object *Object
}

func (o *ObjectBuilder) NewObject(id, class string) *Object {
	obj := NewObject(id, class)
	o.AddObject(obj)
	return obj
}

func (o *ObjectBuilder) AddObject(obj *Object) {
	o.Object = obj
}
