package model

//
// PointerProperty represents a class member which points to a single class instance.
// Unlike RelativeProperty the referenced instance does not know about the reference.
//
type PointerProperty struct {
	id   StringId
	name string
	cls  StringId
}

func NewPointerProperty(id StringId, name string, cls StringId) *PointerProperty {
	return &PointerProperty{id, name, cls}
}

//
// Id returns the unique id of this property.
// It is usually derived from the property's name.
//
func (this *PointerProperty) Id() StringId {
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
func (this *PointerProperty) Class() StringId {
	return this.cls
}
