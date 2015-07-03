package compiler

import (
	M "github.com/ionous/sashimi/model"
)

//
// PendingInstance records all of the classes which use an instance
//
type PendingInstance struct {
	id       M.StringId
	name     string
	longName string
	classes  ClassReferences
}
