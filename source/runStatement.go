package source

import (
	E "github.com/ionous/sashimi/event"
	G "github.com/ionous/sashimi/game"
)

//
// holds event callbacks
type RunStatement struct {
	fields RunFields
	source Code
}

type RunFields struct {
	Owner    string
	Action   string
	Callback G.Callback
	Phase    E.Phase
}

func (this RunStatement) Fields() RunFields {
	return this.fields
}

func (this RunStatement) Source() Code {
	return this.source
}
