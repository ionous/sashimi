package source

//
// An AssertionStatement declares the existence of a class or instance.
//
type AssertionStatement struct {
	fields AssertionFields
	source Code
}

type AssertionFields struct {
	Owner   string // base type or class
	Called  string // name of reference being asserted into existence
	Options        // ex. called
}

//
func (ts AssertionStatement) Fields() AssertionFields {
	return ts.fields
}

//
func (ts AssertionStatement) Source() Code {
	return ts.source
}
