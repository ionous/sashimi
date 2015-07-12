package source

//
type ChoiceStatement struct {
	fields ChoiceFields
	source Code
}

type ChoiceFields struct {
	Owner  string
	Choice string
}

//
func (ts ChoiceStatement) Fields() ChoiceFields {
	return ts.fields
}

//
func (ts ChoiceStatement) Source() Code {
	return ts.source
}
