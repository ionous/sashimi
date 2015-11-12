package model

import "github.com/ionous/sashimi/util/ident"

type NumProperty struct {
	Id     ident.Id `json:"id"`   // property id
	Name   string   `json:"name"` // property name
	IsMany bool     `json:"many"`
}

func (num NumProperty) GetId() ident.Id {
	return num.Id
}

func (num NumProperty) GetName() string {
	return num.Name
}

func (num NumProperty) GetZero(_ ConstraintSet) (ret interface{}) {
	if !num.IsMany {
		ret = float32(0)
	} else {
		ret = []interface{}{}
	}
	return
}
