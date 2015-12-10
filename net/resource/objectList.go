package resource

// By default the object list does not have storage.
// An empty list acts as a "null device".
type ObjectList struct {
	doc *Document
}

// Create a new object list with backing storage
func NewObjectList() ObjectList {
	doc := Document{}
	return ObjectList{&doc}
}

// Add an object identifier to the list of document objects.
// Returns a builder to turn the identifier into a full object.
func (l ObjectList) NewObject(id, class string) *Object {
	obj := NewObject(id, class)
	l.AddObject(obj)
	return obj
}

func (l ObjectList) AddObject(obj *Object) {
	if l.doc != nil {
		l.doc.Included = append(l.doc.Included, obj)
	}
}

// Return the array of added objects, if any.
func (l ObjectList) Objects() (ret []*Object) {
	if l.doc != nil {
		ret = l.doc.Included
	}
	return ret
}
