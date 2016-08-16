package model

import (
	"github.com/ionous/sashimi/util/ident"
	"github.com/ionous/sashimi/util/sbuf"
)

type InstanceModel struct {
	Id     ident.Id `json:"id"`
	Class  ident.Id `json:"type"`
	Name   string   `json:"name"`
	Values Values   `json:"values,omitempty"`
}

type Values map[ident.Id]Value

// Enums are stored as int;
// Numbers as float64;
// Pointers as ident.Id;
// Text as string.
type Value interface{}

func (n *InstanceModel) String() string {
	return sbuf.New().S(string(n.Id)).R('(').S(string(n.Class)).R(')').String()
}
func (n *InstanceModel) GetId() ident.Id {
	return n.Id
}
func (n *InstanceModel) GetParentClass() ident.Id {
	return n.Class
}
func (n *InstanceModel) GetOriginalName() string {
	return n.Name
}
