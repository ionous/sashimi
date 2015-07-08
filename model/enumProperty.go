package model

import "github.com/ionous/sashimi/util/ident"

type EnumProperty struct {
	id   ident.Id
	name string
	Enumeration
}

func NewEnumProperty(id ident.Id, name string, src Enumeration) *EnumProperty {
	return &EnumProperty{id, name, src}
}

func (this *EnumProperty) Id() ident.Id {
	return this.id
}

func (this *EnumProperty) Name() string {
	return this.name
}
