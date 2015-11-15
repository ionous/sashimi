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
	List  bool
}

//
func (ts PropertyStatement) Fields() PropertyFields {
	return ts.fields
}

//
func (ts PropertyStatement) Source() Code {
	return ts.source
}
