package source

type ActionStatement struct {
	fields ActionAssertionFields
	source Code
}

func (this ActionStatement) Source() Code {
	return this.source
}

func (this ActionStatement) Fields() ActionAssertionFields {
	return this.fields
}

type ActionAssertionFields struct {
	ActionName string
	EventName  string
	Source     string
	Target     string
	Context    string
}
