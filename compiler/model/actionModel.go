package model

import "github.com/ionous/sashimi/util/ident"

type ActionModel struct {
	Id             ident.Id   `json:"id"`
	Name           string     `json:"name"`    // Original Name
	EventId        ident.Id   `json:"event"`   // Related Event
	NounTypes      []ident.Id `json:"nouns"`   // Classes
	DefaultActions []ident.Id `json:"actions"` // Callbacks
}

func (a ActionModel) String() string {
	return a.Name
}
