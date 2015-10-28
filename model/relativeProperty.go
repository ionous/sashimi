package model

import "github.com/ionous/sashimi/util/ident"

type RelativeProperty struct {
	Id       ident.Id `json:"id"`       // property id
	Name     string   `json:"name"`     // property name
	Class    ident.Id `json:"class"`    // property id
	Relates  ident.Id `json:"relates"`  // other class id
	Relation ident.Id `json:"relation"` // relation id
	IsRev    bool     `json:"rev"`
	IsMany   bool     `json:"many"`
}

func (rel RelativeProperty) GetId() ident.Id {
	return rel.Id
}

func (rel RelativeProperty) GetName() string {
	return rel.Name
}

func (rel RelativeProperty) GetZero(_ ConstraintSet) interface{} {
	return nil
}
