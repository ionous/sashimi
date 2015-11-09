package memory

import "github.com/ionous/sashimi/util/ident"

//
type GenericValue interface{}

//
type ObjectValue interface {
	GetValue(obj, field ident.Id) (GenericValue, bool)
	SetValue(obj, field ident.Id, v GenericValue) error
}
