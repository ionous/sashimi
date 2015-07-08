package model

import "github.com/ionous/sashimi/util/ident"

type IProperty interface {
	Id() ident.Id
	Name() string
}

type PropertySet map[ident.Id]IProperty
