package source

//
type KeyValueStatement struct {
	fields KeyValueFields
	source Code
}

type KeyValueFields struct {
	Owner string
	Key   string
	Value interface{}
}

//
func (this KeyValueStatement) Fields() KeyValueFields {
	return this.fields
}

//
func (this KeyValueStatement) Source() Code {
	return this.source
}
