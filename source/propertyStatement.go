package source

//
type PropertyStatement struct {
	fields PropertyFields
	source Code
}

type PropertyFields struct {
	Class string
	Name  string
	Kind  string
}

//
func (this PropertyStatement) Fields() PropertyFields {
	return this.fields
}

//
func (this PropertyStatement) Source() Code {
	return this.source
}
