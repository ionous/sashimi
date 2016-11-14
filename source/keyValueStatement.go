package source

import (
	. "github.com/ionous/sashimi/source/types"
)

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
func (ts KeyValueStatement) Fields() KeyValueFields {
	return ts.fields
}

//
func (ts KeyValueStatement) Source() Code {
	return ts.source
}
