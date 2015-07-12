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

func (enum *EnumProperty) Id() ident.Id {
	return enum.id
}

func (enum *EnumProperty) Name() string {
	return enum.name
}

func (enum *EnumProperty) Zero(cons ConstraintSet) interface{} {
	index := 1 // enum indices are 1-based.
	if cons, ok := cons.ConstraintById(enum.id); ok {
		if cons, ok := cons.(*EnumConstraint); ok {
			index = cons.BestIndex()
		}
	}
	if _, e := enum.IndexToChoice(index); e != nil {
		panic(e)
	}
	return index
}
