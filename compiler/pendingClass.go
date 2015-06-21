package compiler

import (
	M "github.com/ionous/sashimi/model"
	S "github.com/ionous/sashimi/source"
)

type PendingClass struct {
	owner     *ClassFactory
	parent    *PendingClass
	id        M.StringId
	name      string
	singular  string
	names     NameScope
	props     PendingProperties
	rules     PendingRules
	relatives PendingRelatives
}

type PendingRelativeEntry struct {
	src S.Code
	rel PendingRelative
}
type PendingRelatives map[M.StringId]PendingRelativeEntry

//
func (this *PendingClass) String() string {
	return this.name
}
