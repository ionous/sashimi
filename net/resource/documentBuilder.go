package resource

import "fmt"

// DocumentBuilder provides an api for creating json-api document structures.
type DocumentBuilder struct {
	doc *Document
}

//
// NewDocumentBuilder creates an api-object to hepl populate the passed document.
//
func NewDocumentBuilder(doc *Document) DocumentBuilder {
	return DocumentBuilder{doc}
}

//
// NewObject sets the document data to the passed object identifier.
// Returns an object builder to populate data about the identified object.
//
func (build DocumentBuilder) NewObject(id, class string) *Object {
	obj := NewObject(id, class)
	build.AddObject(obj)
	return obj
}

func (build DocumentBuilder) AddObject(obj *Object) {
	if build.doc.Data == nil {
		build.doc.Data = obj
	} else {
		build.AddError(fmt.Errorf("document object specified multiple times."))
	}
}

//
// NewObjects returns a builder to add an array of objects ( or object identifiers ) to the document.
//
func (build DocumentBuilder) NewObjects() ObjectsBuilder {
	if build.doc.Data == nil {
		// a blank placeholder for ObjectsBuilder.NewObject
		build.doc.Data = []Object{}
	} else {
		build.AddError(fmt.Errorf("document objects specified multiple times."))
	}
	return ObjectsBuilder{build}
}

//
// Sets the document data to an existing array of objects or object identifiers.
// Returns a builder to add objects to that array.
//
func (build DocumentBuilder) Sets(objects ObjectList) DocumentBuilder {
	if build.doc.Data == nil {
		// a blank placeholder for ObjectsBuilder.NewObject
		build.doc.Data = objects.doc.Included
	} else {
		build.AddError(fmt.Errorf("document objects specified multiple times."))
	}
	return build
}

//
// SetMeta to add a key-value the document's metadata.
//
func (build DocumentBuilder) SetMeta(key string, value interface{}) DocumentBuilder {
	if build.doc.Meta == nil {
		build.doc.Meta = Dict{}
	}
	build.doc.Meta[key] = value
	return build
}

//
// SetLink to add the named link to the document's list of links.
//
func (build DocumentBuilder) SetLink(key string, link Link) DocumentBuilder {
	if build.doc.Links == nil {
		build.doc.Links = make(Links)
	}
	build.doc.Links[key] = link
	return build
}

//
// NewIncludes returns a builder which can add objects (or object identifiers) to a compound document.
//
func (build DocumentBuilder) NewIncludes() ObjectList {
	if build.doc.Included == nil {
		build.doc.Included = []*Object{}
	} else {
		build.AddError(fmt.Errorf("document objects included multiple times."))
	}
	return ObjectList{build.doc}
}

//
// SetIncluded sets the document's compound/include data to the passed object list.
// ( To comply with jsonapi, objects in the list should be referenced by the primary object. )
//
func (build DocumentBuilder) SetIncluded(objects ObjectList) DocumentBuilder {
	if build.doc.Included == nil {
		build.doc.Included = objects.Objects()
	} else {
		build.AddError(fmt.Errorf("document objects included multiple times."))
	}
	return build
}

//
// AddError appends an error to the list of errors included by this document.
// NOTE: if ever needed, could return or take error builder
// to which the other bits of the jsonapi error structure could be added
//
func (build DocumentBuilder) AddError(err error) DocumentBuilder {
	build.doc.Errors = append(build.doc.Errors, Error{Code: err.Error()})
	return build
}
