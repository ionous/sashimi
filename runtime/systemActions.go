package runtime

import (
	"github.com/ionous/sashimi/util/ident"
)

// SystemActions are run after all other actions have taken place
type SystemActions map[ident.Id][]SystemCallback

type SystemCallback func(obj []*GameObject)

func (a SystemActions) Capture(act ident.Id, cb SystemCallback) {
	arr := a[act]
	arr = append(arr, cb)
	a[act] = arr
}

func (a SystemActions) Run(act ident.Id, objects []*GameObject) {
	if arr, ok := a[act]; ok {
		for _, cb := range arr {
			cb(objects)
		}
	}
}
