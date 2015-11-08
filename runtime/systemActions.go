package runtime

import (
	"fmt"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/util/ident"
)

// SystemActions are run after all other actions have taken place
type SystemActionMap map[ident.Id][]SystemCallback
type SystemActions struct {
	mdl     api.Model
	actions SystemActionMap
}

type SystemCallback func(obj []*GameObject)

func (a SystemActions) Capture(act ident.Id, cb SystemCallback) (err error) {
	if _, ok := a.mdl.GetAction(act); !ok {
		err = fmt.Errorf("couldnt find action %s", act)
	} else {
		arr := a.actions[act]
		arr = append(arr, cb)
		a.actions[act] = arr
	}
	return err
}

func (a SystemActions) Run(act ident.Id, objects []*GameObject) {
	if arr, ok := a.actions[act]; ok {
		for _, cb := range arr {
			cb(objects)
		}
	}
}
