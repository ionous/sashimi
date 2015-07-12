package source

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
