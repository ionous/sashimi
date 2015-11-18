package xmodel

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
type ConstraintSet struct {
	Parent *ConstraintSet
	Map    ConstraintMap
}
type ConstraintMap map[ident.Id]IConstrain

func (cons ConstraintSet) Len() int {
	return len(cons.Map)
}
func (cons ConstraintSet) ConstraintById(id ident.Id) (ret IConstrain, okay bool) {
	if c, ok := cons.Map[id]; ok {
		ret, okay = c, ok
	} else if cons.Parent != nil {
		ret, okay = cons.Parent.ConstraintById(id)
	}
	return ret, okay
}

//
type UnknownConstraintError struct {
	prop       IProperty
	constraint IConstrain
}

func (this UnknownConstraintError) Error() string {
	return fmt.Sprintf("unknown constraint %v for %v", this.constraint, this.prop)
}
