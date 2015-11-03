package model

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

func (enum EnumProperty) GetZero(cons ConstraintSet) interface{} {
	index := 1 // enum indices are 1-based.
	if cons, ok := cons.ConstraintById(enum.Id); ok {
		if cons, ok := cons.(*EnumConstraint); ok {
			index = cons.BestIndex()
		}
	}
	if _, e := enum.IndexToChoice(index); e != nil {
		panic(e)
	}
	return index
}
