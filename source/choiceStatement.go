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
func (this ChoiceStatement) Fields() ChoiceFields {
	return this.fields
}

//
func (this ChoiceStatement) Source() Code {
	return this.source
}
