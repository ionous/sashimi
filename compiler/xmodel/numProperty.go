package xmodel

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
