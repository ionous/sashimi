package memory

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/util/ident"
)

type eventInfo struct {
	mdl *MemoryModel
	*M.EventInfo
}

func (e eventInfo) GetId() ident.Id {
	return e.Id
}

func (e eventInfo) GetEventName() string {
	return e.EventName
}

func (e eventInfo) GetAction() (ret api.Action) {
	if a, ok := e.mdl.GetAction(e.ActionId); !ok {
		panic(fmt.Sprintf("internal error, no action found for event %s", e.ActionId))
	} else {
		ret = a
	}
	return
}

func (e eventInfo) GetListeners(capture bool) (ret api.Listeners, okay bool) {
	var callbacks EventCallbacks
	if !capture {
		callbacks = e.mdl.bubble
	} else {
		callbacks = e.mdl.capture
	}
	if cbs, ok := callbacks[e.Id]; !ok {
		ret = api.NoListeners{}
	} else {
		ret = listenersInfo{e.mdl, e.EventInfo, cbs, capture}
		okay = true
	}
	return
}
