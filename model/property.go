package model

import "github.com/ionous/sashimi/util/ident"

//
// IProperty represents a sashimi type.
//
type IProperty interface {
	Id() ident.Id
	Name() string
	// note: the determination of zero value is not possible in a purely generic way.
	// a property's zero value requires the constraints provided by its class.
	Zero(ConstraintSet) interface{}
}

type PropertySet map[ident.Id]IProperty
