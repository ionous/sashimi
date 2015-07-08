package model

import "github.com/ionous/sashimi/util/ident"

type NumProperty struct {
	id   ident.Id
	name string
}

func NewNumProperty(id ident.Id, name string) *NumProperty {
	return &NumProperty{id, name}
}

func (this *NumProperty) Id() ident.Id {
	return this.id
}

func (this *NumProperty) Name() string {
	return this.name
}
