package source

//
type AliasStatement struct {
	fields AliasFields
	source Code
}

type AliasFields struct {
	Key     string
	Phrases []string
}

//
func (ts AliasStatement) Fields() AliasFields {
	return ts.fields
	// since we can take the address of a field and write to that address
	// even key isnt truly read-only, so why worry about arrays?
	// AliasFields{ts.fields.Key, copy(ts.fields.Phrases)}
}

//
func (ts AliasStatement) Source() Code {
	return ts.source
}
