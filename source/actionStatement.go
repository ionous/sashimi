package source

import (
	"github.com/ionous/sashimi/source/types"
	"github.com/ionous/sashimi/util/sbuf"
)

type ActionStatement struct {
	fields ActionAssertionFields
	source types.Code
}

func (ts ActionStatement) Source() types.Code {
	return ts.source
}

func (ts ActionStatement) Fields() ActionAssertionFields {
	return ts.fields
}

type ActionAssertionFields struct {
	Action  types.NamedAction
	Event   types.NamedEvent
	Source  string
	Target  types.NamedClass
	Context types.NamedClass
}

func (f ActionAssertionFields) String() string {
	return sbuf.New("Action:", f.Action,
		",Event:", f.Event,
		",Source:", f.Source,
		",Target:", f.Target,
		",Context:", f.Context).String()
}
