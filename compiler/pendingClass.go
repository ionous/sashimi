package compiler

import (
	M "github.com/ionous/sashimi/model"
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

type PendingRelatives map[M.StringId]PendingRelative

//
func (this *PendingClass) String() string {
	return this.name
}
