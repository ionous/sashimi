package memory

import (
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/util/ident"
)

type eventInfo struct {
	mdl *MemoryModel
	*M.EventModel
}

func (e eventInfo) GetId() ident.Id {
	return e.Id
}

func (e eventInfo) GetEventName() string {
	return e.Name
}

func (e eventInfo) GetListeners(capture bool) (ret api.Listeners, okay bool) {
	var callbacks M.EventModelCallbacks
	if !capture {
		callbacks = e.Bubble
	} else {
		callbacks = e.Capture
	}
	if len(callbacks) == 0 {
		ret = api.NoListeners{}
	} else {
		ret = listenersInfo{e.mdl, e.EventModel, callbacks, capture}
		okay = true
	}
	return
}
