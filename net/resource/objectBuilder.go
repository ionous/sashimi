package resource

//
// DocumentBuilder, ObjectsBuilder, and ObjectList satisfy this interface,
// which allows for either one object or multiple objects to be added depending on the context.
//
type IBuildObjects interface {
	NewObject(id, class string) *Object
}

type ObjectBuilder struct {
	Object *Object
}

func (this *ObjectBuilder) NewObject(id, class string) *Object {
	this.Object = NewObject(id, class)
	return this.Object
}
