package model

import "github.com/ionous/sashimi/util/ident"

type PropertyType string

const (
	NumProperty      PropertyType = "NumProperty"
	TextProperty                  = "TextProperty"
	PointerProperty               = "PointerProperty"
	EnumProperty                  = "EnumProperty"
	RelativeProperty              = "RelativeProperty"
)

type PropertyModel struct {
	Id   ident.Id     `json:"id"`
	Type PropertyType `json:"type"`
	Name string       `json:"name"`
	// pointer, relation
	Relates ident.Id `json:"relates,omitempty"`
	// relation
	Relation ident.Id `json:"relation,omitempty"` // relation id
	// num, text, pointer
	IsMany bool `json:"many,omitempty"`
	// relation
	IsRev bool `json:"rev,omitempty"`
}
