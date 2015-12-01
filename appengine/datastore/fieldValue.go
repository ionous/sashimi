package datastore

import (
	"github.com/ionous/sashimi/meta"
)

// FieldValue knows how to serialize an value into the datastore.
type FieldValue struct {
	Type  meta.PropertyType
	Value interface{}
}

// Encode makes a datastore-able value.
func (iv FieldValue) Encode() (ret interface{}, err error) {
	return Encode(iv.Type, iv.Value)
}
