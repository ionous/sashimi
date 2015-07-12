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

func (rel *RelativeProperty) Fields() RelativeFields {
	return rel.fields
}

func (rel *RelativeProperty) Class() ident.Id {
	return rel.fields.Class
}

func (rel *RelativeProperty) Id() ident.Id {
	return rel.fields.Id
}

func (rel *RelativeProperty) Name() string {
	return rel.fields.Name
}

// id of the relation table
func (rel *RelativeProperty) Relation() ident.Id {
	return rel.fields.Relation
}

// id of the other class rel property inolves
func (rel *RelativeProperty) Relates() ident.Id {
	return rel.fields.Relates
}

// in the case where a one to many relation involves a class and itself,
// distingushes which side of the relation propery is the primary.
func (rel *RelativeProperty) IsRev() bool {
	return rel.fields.IsRev
}

// in the case where a one to many relation involves a class and itself,
// distingushes which side of the relation propery is the many.
func (rel *RelativeProperty) ToMany() bool {
	return rel.fields.ToMany
}

func (rel *RelativeProperty) Zero(_ ConstraintSet) interface{} {
	return nil
}
