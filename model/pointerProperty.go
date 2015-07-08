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
// Id returns the unique id of this property.
// It is usually derived from the property's name.
//
func (this *PointerProperty) Id() ident.Id {
	return this.id
}

//
// Name returns the property's appelation as specified by the author.
//
func (this *PointerProperty) Name() string {
	return this.name
}

//
// Class returns the id of the kind of class which instances with this property point to.
//
func (this *PointerProperty) Class() ident.Id {
	return this.cls
}
