package model

import "github.com/ionous/sashimi/util/ident"

type EventModel struct {
	Id      ident.Id            `json:"id"`
	Name    string              `json:"name"`
	Capture EventModelCallbacks `json:"capture,omitempty"`
	Bubble  EventModelCallbacks `json:"bubble,omitempty"`
}

func (e EventModel) String() string {
	return e.Name
}

type EventModelCallbacks []ListenerModel
