package model

import (
	"fmt"
	"github.com/ionous/sashimi/util/ident"
)

type IConstrain interface {
	Always(value interface{}) error
	Usually(value interface{}) error
	Seldom(value interface{}) error
	Exclude(value interface{}) error
	Copy() IConstrain
}

//
//
//
type ConstraintSet struct {
	parent      *ConstraintSet
	constraints ConstraintMap
}
type ConstraintMap map[ident.Id]IConstrain

func (cons ConstraintSet) ConstraintById(id ident.Id) (ret IConstrain, okay bool) {
	if c, ok := cons.constraints[id]; ok {
		ret, okay = c, ok
	} else if cons.parent != nil {
		ret, okay = cons.parent.ConstraintById(id)
	}
	return ret, okay
}

//
//
//
type UnknownConstraintError struct {
	prop       IProperty
	constraint IConstrain
}

func (this UnknownConstraintError) Error() string {
	return fmt.Sprintf("unknown constraint %v for %v", this.constraint, this.prop)
}
