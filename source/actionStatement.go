package source

import "github.com/ionous/sashimi/util/sbuf"

type ActionStatement struct {
	fields ActionAssertionFields
	source Code
}

func (ts ActionStatement) Source() Code {
	return ts.source
}

func (ts ActionStatement) Fields() ActionAssertionFields {
	return ts.fields
}

type ActionAssertionFields struct {
	Action  string
	Event   string
	Source  string
	Target  string
	Context string
}

func (f ActionAssertionFields) String() string {
	return sbuf.New("Action:", f.Action,
		",Event:", f.Event,
		",Source:", f.Source,
		",Target:", f.Target,
		",Context:", f.Context).String()
}
