package xmodel

import (
	"fmt"
	"github.com/ionous/sashimi/util/errutil"
	"github.com/ionous/sashimi/util/ident"
)

type ActionMap map[ident.Id]*ActionInfo

// FIX: i kind of think all maps with the string ids should have this
func (this ActionMap) FindActionByName(name string) (*ActionInfo, bool) {
	id := MakeStringId(name)
	act, ok := this[id]
	return act, ok
}

//
func ActionNotFound(str string) error {
	return errutil.Func(func() string {
		return fmt.Sprintf("unknown action requested %s", str)
	})
}
