package model

import "github.com/ionous/sashimi/util/ident"

type TextProperty struct {
	Id   ident.Id `json:"id"`   // property id
	Name string   `json:"name"` // property name
}

func (text TextProperty) GetId() ident.Id {
	return text.Id
}

func (text TextProperty) GetName() string {
	return text.Name
}

func (text TextProperty) GetZero(_ ConstraintSet) interface{} {
	return ""
}
