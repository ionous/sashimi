package model

import (
	"fmt"
	"github.com/ionous/sashimi/util/errutil"
	"github.com/ionous/sashimi/util/ident"
)

type InstanceMap map[ident.Id]*InstanceInfo

// helper to generate an escaped string and an error,
// FIX: to generize the name to string search with distance, you will (probably) have to walk all strings
// a generic filter function with which takes a distance visitor might be the best you can do in go.
func (this InstanceMap) FindInstance(name string) (*InstanceInfo, bool) {
	safe := MakeStringId(name)
	ret, okay := this[safe]
	return ret, okay
}

func InstanceNotFound(name string) error {
	return errutil.Func(func() string {
		return fmt.Sprintf("instance not found `%s`", name)
	})
}
