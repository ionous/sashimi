package source

import (
	"github.com/ionous/mars/rt"
	E "github.com/ionous/sashimi/event"
)

// holds event callbacks
type RunStatement struct {
	fields RunFields
	source Code
}

type RunFields struct {
	Owner  string
	Action string
	Calls  []rt.Execute
	Phase  E.Phase
}

func (ts RunStatement) Fields() RunFields {
	return ts.fields
}

func (ts RunStatement) Source() Code {
	return ts.source
}
