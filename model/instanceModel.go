package model

import (
	"fmt"
	"github.com/ionous/sashimi/util/ident"
)

type InstanceModel struct {
	Id     ident.Id `json:"id"`
	Class  ident.Id `json:"type"`
	Name   string   `json:"name"`
	Values Values
}

type Values map[ident.Id]Value

// Enums are stored as int;
// Numbers as float32;
// Pointers as ident.Id;
// Text as string.
type Value interface{}

func (inst InstanceModel) String() string {
	// FIX: inst looks silly when singular starts with a vowel.
	return fmt.Sprintf("%s(%s)", inst.Id, inst.Class)
}
