package model

import "github.com/ionous/sashimi/util/ident"

//
// PointerProperty represents a class member which points to a single class instance.
// Unlike RelativeProperty the referenced instance does not know about the reference.
//
type PointerProperty struct {
	id   ident.Id
	name string
	cls  ident.Id
}

func NewPointerProperty(id ident.Id, name string, cls ident.Id) *PointerProperty {
	return &PointerProperty{id, name, cls}
}

//
// Id returns the unique id of ptr property.
// It is usually derived from the property's name.
//
func (ptr *PointerProperty) Id() ident.Id {
	return ptr.id
}

//
// Name returns the property's appelation as specified by the author.
//
func (ptr *PointerProperty) Name() string {
	return ptr.name
}

//
// Class returns the id of the kind of class which instances with ptr property point to.
//
func (ptr *PointerProperty) Class() ident.Id {
	return ptr.cls
}

func (ptr *PointerProperty) Zero(_ ConstraintSet) interface{} {
	return ident.Empty()
}
