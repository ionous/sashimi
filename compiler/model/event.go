package model

import "github.com/ionous/sashimi/util/ident"

type EventModel struct {
	Id              ident.Id
	Name            string // Original Name
	Capture, Bubble EventModelCallbacks
}

func (e EventModel) String() string {
	return e.Name
}

type EventModelCallbacks []ListenerModel
