package metal

import (
	"github.com/ionous/sashimi/util/ident"
)

//
type GenericValue interface{}

// another way, potentially, of handling this would be to change the model in memory; and use property watchers to record to the db.
type ObjectValue interface {
	GetValue(obj, field ident.Id) (GenericValue, bool)
	SetValue(obj, field ident.Id, v GenericValue) error
}
