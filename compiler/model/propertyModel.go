package model

import "github.com/ionous/sashimi/util/ident"

type PropertyType string

const (
	NumProperty     PropertyType = "NumProperty"
	TextProperty                 = "TextProperty"
	PointerProperty              = "PointerProperty"
	EnumProperty                 = "EnumProperty"
)

type PropertyModel struct {
	Id   ident.Id     `json:"id"`
	Type PropertyType `json:"type"`
	Name string       `json:"name"`
	// pointer, relation
	Relates ident.Id `json:"relates,omitempty"`
	// relation
	Relation ident.Id `json:"relation,omitempty"`
	// num, text, pointer, relation; but, not enum
	IsMany bool `json:"many,omitempty"`
}
