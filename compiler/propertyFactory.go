package compiler

import (
	M "github.com/ionous/sashimi/model"
	S "github.com/ionous/sashimi/source"
)

//
// Primitive types which can be specified by an author.
//
const (
	TextType = "text"
	NumType  = "num"
)

//
// The in-progress properties of a single class
//
type PendingProperties map[M.StringId]M.IProperty

type PendingRelativeEntry struct {
	src S.Code
	rel PendingRelative
}
type PendingRelatives map[M.StringId]PendingRelativeEntry
