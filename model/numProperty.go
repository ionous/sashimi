package model

import "github.com/ionous/sashimi/util/ident"

type NumProperty struct {
	id   ident.Id
	name string
}

func NewNumProperty(id ident.Id, name string) *NumProperty {
	return &NumProperty{id, name}
}

func (num *NumProperty) Id() ident.Id {
	return num.id
}

func (num *NumProperty) Name() string {
	return num.name
}

func (num *NumProperty) Zero(_ ConstraintSet) interface{} {
	return float32(0)
}
