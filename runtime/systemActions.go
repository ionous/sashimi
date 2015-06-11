package runtime

import (
	M "github.com/ionous/sashimi/model"
)

func NewSystemActions() SystemActions {
	return SystemActions{make(map[string][]SystemCallback)}
}

type SystemActions struct {
	actions map[string][]SystemCallback
}

type SystemCallback func(action *M.ActionInfo, obj []*GameObject)

func (this *SystemActions) Capture(event string, cb SystemCallback) *SystemActions {
	arr := this.actions[event]
	arr = append(arr, cb)
	this.actions[event] = arr
	return this
}

func (this *SystemActions) Run(action *M.ActionInfo, obj []*GameObject) {
	if arr, ok := this.actions[action.Event()]; ok {
		for _, cb := range arr {
			cb(action, obj)
		}
	}
}
