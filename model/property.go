package model

type IProperty interface {
	Id() StringId
	Name() string
}

type PropertySet map[StringId]IProperty
