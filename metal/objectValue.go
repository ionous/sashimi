package metal

import (
	"github.com/ionous/sashimi/util/ident"
)

// for mocks
type ObjectValue interface {
	GetValue(cls, field ident.Id) (interface{}, bool)
	SetValue(cls, field ident.Id, v interface{}) error
}
