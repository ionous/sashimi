package model

import "fmt"
import "github.com/ionous/sashimi/util/ident"

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
type ConstraintSet map[ident.Id]IConstrain

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
