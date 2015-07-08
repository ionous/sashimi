package model

import "fmt"
import "github.com/ionous/sashimi/util/ident"

type ActionMap map[ident.Id]*ActionInfo

// FIX: i kind of think all maps with the string ids should have this
func (this ActionMap) FindActionByName(name string) (*ActionInfo, bool) {
	id := MakeStringId(name)
	act, ok := this[id]
	return act, ok
}

//
type ActionNotFound string

//
func (this ActionNotFound) Error() string {
	return fmt.Sprintf("unknown action requested %s", this)
}
