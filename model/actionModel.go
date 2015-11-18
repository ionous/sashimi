package model

import "github.com/ionous/sashimi/util/ident"

type ActionModel struct {
	Id             ident.Id
	Name           string     // Original Name
	EventId        ident.Id   // Related Event
	NounTypes      []ident.Id // Classes
	DefaultActions []ident.Id // Callbacks
}

func (a ActionModel) String() string {
	return a.Name
}
