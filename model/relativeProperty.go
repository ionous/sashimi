package model

import "github.com/ionous/sashimi/util/ident"

// NOTE: Class, IsMany. Relates are redudent/congruant with information in the relation
type RelativeProperty struct {
	Id       ident.Id `json:"id"`       // property id
	Name     string   `json:"name"`     // property name
	Relation ident.Id `json:"relation"` // relation id
	IsRev    bool     `json:"rev"`      // when true, this property describes relation.Dest
	//
	Relates ident.Id `json:"relates"` // other class id
	IsMany  bool     `json:"many"`    // when true, refers to many object
}

func (rel RelativeProperty) GetId() ident.Id {
	return rel.Id
}

func (rel RelativeProperty) GetName() string {
	return rel.Name
}

func (rel RelativeProperty) GetZero(_ ConstraintSet) interface{} {
	panic("relative properties dont have values")
}
