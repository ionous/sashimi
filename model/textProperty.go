package model

import "github.com/ionous/sashimi/util/ident"

type TextProperty struct {
	id   ident.Id
	name string
}

// FIX?  it might be nicer if the model had a builder with an interface for all such new(s)
func NewTextProperty(id ident.Id, name string) *TextProperty {
	return &TextProperty{id, name}
}

func (this *TextProperty) Id() ident.Id {
	return this.id
}

func (this *TextProperty) Name() string {
	return this.name
}
