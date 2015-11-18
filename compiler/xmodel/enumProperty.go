package xmodel

import "github.com/ionous/sashimi/util/ident"

type EnumProperty struct {
	Id   ident.Id `json:"id"`   // property id
	Name string   `json:"name"` // property name
	Enumeration
}

func (enum EnumProperty) GetId() ident.Id {
	return enum.Id
}

func (enum EnumProperty) GetName() string {
	return enum.Name
}
