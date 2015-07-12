package resource

import "fmt"

type DocumentBuilder struct {
	doc *Document
}

func NewDocumentBuilder(doc *Document) DocumentBuilder {
	return DocumentBuilder{doc}
}

//
// Set the document data to the passed object identifier.
// Return a builder to turn the identifier into a full object.
//
func (this DocumentBuilder) NewObject(id, class string) *Object {
	obj := NewObject(id, class)
	if this.doc.Data == nil {
		this.doc.Data = obj
	} else {
		this.AddError(fmt.Errorf("document object specified multiple times."))
	}
	return obj
}

//
// Set the document data to an array of objects or object identifiers.
// Return a builder to add objects to that array.
//
func (this DocumentBuilder) NewObjects() ObjectsBuilder {
	if this.doc.Data == nil {
		// a blank placeholder for ObjectsBuilder.NewObject
		this.doc.Data = []Object{}
	} else {
		this.AddError(fmt.Errorf("document objects specified multiple times."))
	}
	return ObjectsBuilder{this}
}

//
// Set the document data to an array of objects or object identifiers.
// Return a builder to add objects to that array.
//
func (this DocumentBuilder) Sets(objects ObjectList) DocumentBuilder {
	if this.doc.Data == nil {
		// a blank placeholder for ObjectsBuilder.NewObject
		this.doc.Data = objects.doc.Included
	} else {
		this.AddError(fmt.Errorf("document objects specified multiple times."))
	}
	return this
}

//
// Add metadata to document.
//
func (this DocumentBuilder) SetMeta(key string, value interface{}) DocumentBuilder {
	if this.doc.Meta == nil {
		this.doc.Meta = Dict{}
	}
	this.doc.Meta[key] = value
	return this
}

//
// Add a link to the document.
//
func (this DocumentBuilder) SetLink(key string, link Link) DocumentBuilder {
	if this.doc.Links == nil {
		this.doc.Links = make(Links)
	}
	this.doc.Links[key] = link
	return this
}

//
// Set the document data to an array of objects or object identifiers.
// Return a builder to add objects to that array.
//
func (this DocumentBuilder) NewIncludes() ObjectList {
	if this.doc.Included == nil {
		this.doc.Included = []*Object{}
	} else {
		this.AddError(fmt.Errorf("document objects included multiple times."))
	}
	return ObjectList{this.doc}
}

//
// Add additional objects to the document.
// ( They should be referenced by the primary object in someway. )
//
func (this DocumentBuilder) SetIncluded(objects ObjectList) DocumentBuilder {
	if this.doc.Included == nil {
		this.doc.Included = objects.Objects()
	} else {
		this.AddError(fmt.Errorf("document objects included multiple times."))
	}
	return this
}

//
// NOTE: if ever needed, could return or take error builder
// to which the other bits of the jsonapi error structure could be added
//
func (this DocumentBuilder) AddError(err error) DocumentBuilder {
	this.doc.Errors = append(this.doc.Errors, Error{Code: err.Error()})
	return this
}
