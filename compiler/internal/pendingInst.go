package internal

import (
	"github.com/ionous/sashimi/util/ident"
)

//
// PendingInstance records all of the classes which use an instance
//
type PendingInstance struct {
	id       ident.Id
	name     string
	longName string
	classes  ClassReferences
}
