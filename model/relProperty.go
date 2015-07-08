package model

import "github.com/ionous/sashimi/util/ident"

type RelativeProperty struct {
	fields RelativeFields
}

type RelativeFields struct {
	Class    ident.Id `json:"id"`       // property id
	Id       ident.Id `json:"id"`       // property id
	Name     string   `json:"name"`     // property name
	Relates  ident.Id `json:"relates"`  // other class id
	Relation ident.Id `json:"relation"` // relation id
	IsRev    bool     `json:"rev"`
	ToMany   bool     `json:"many"`
}

func NewRelativeProperty(fields RelativeFields) *RelativeProperty {
	return &RelativeProperty{fields}
}

func (this *RelativeProperty) Class() ident.Id {
	return this.fields.Class
}

func (this *RelativeProperty) Id() ident.Id {
	return this.fields.Id
}

func (this *RelativeProperty) Name() string {
	return this.fields.Name
}

// id of the relation table
func (this *RelativeProperty) Relation() ident.Id {
	return this.fields.Relation
}

// id of the other class this property inolves
func (this *RelativeProperty) Relates() ident.Id {
	return this.fields.Relates
}

// in the case where a one to many relation involves a class and itself,
// distingushes which side of the relation propery is the primary.
func (this *RelativeProperty) IsRev() bool {
	return this.fields.IsRev
}

// in the case where a one to many relation involves a class and itself,
// distingushes which side of the relation propery is the many.
func (this *RelativeProperty) ToMany() bool {
	return this.fields.ToMany
}
