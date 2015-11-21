package datastore

import (
	"github.com/ionous/sashimi/meta"
)

type FieldValue struct {
	Type  meta.PropertyType
	Value interface{}
}

func (iv FieldValue) Encode() (ret interface{}, err error) {
	return Encode(iv.Type, iv.Value)
}
