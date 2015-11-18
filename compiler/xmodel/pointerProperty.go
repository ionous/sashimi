package xmodel

import "github.com/ionous/sashimi/util/ident"

//
// PointerProperty represents a class member which points to a single class instance.
// Unlike RelativeProperty the referenced instance does not know about the reference.
//
type PointerProperty struct {
	Id     ident.Id `json:"id"`    // property id
	Name   string   `json:"name"`  // property name
	Class  ident.Id `json:"class"` // property id
	IsMany bool     `json:"many"`  // property id
}

// Id returns the unique id of ptr property.
// It is usually derived from the property's name.
func (ptr PointerProperty) GetId() ident.Id {
	return ptr.Id
}

// Name returns the property's appelation as specified by the author.
func (ptr PointerProperty) GetName() string {
	return ptr.Name
}
