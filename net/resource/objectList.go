package resource

//
// By default the object list does not have storage.
// An empty list acts as a "null device".
//
type ObjectList struct {
	doc *Document
}

//
// Create a new object list with backing storage
//
func NewObjectList() ObjectList {
	doc := Document{}
	return ObjectList{&doc}
}

//
// Add an object identifier to the list of document objects.
// Return a builder to turn the identifier into a full object.
//
func (this ObjectList) NewObject(id, class string) *Object {
	obj := NewObject(id, class)
	if this.doc != nil {
		this.doc.Included = append(this.doc.Included, obj)
	}
	return obj
}

//
// Return the array of added objects, if any.
//
func (this ObjectList) Objects() (ret []*Object) {
	if this.doc != nil {
		ret = this.doc.Included
	}
	return ret
}
