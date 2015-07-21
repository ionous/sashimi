package source

import (
	E "github.com/ionous/sashimi/event"
	G "github.com/ionous/sashimi/game"
)

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

func (ts RunStatement) Fields() RunFields {
	return ts.fields
}

func (ts RunStatement) Source() Code {
	return ts.source
}
