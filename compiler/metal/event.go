package metal

import (
	M "github.com/ionous/sashimi/compiler/model"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
)

type eventInfo struct {
	mdl *Metal
	*M.EventModel
}

func (e eventInfo) GetId() ident.Id {
	return e.Id
}

func (e eventInfo) GetEventName() string {
	return e.Name
}

func (e eventInfo) GetListeners(capture bool) (ret meta.Listeners, okay bool) {
	var callbacks M.EventModelCallbacks
	if !capture {
		callbacks = e.Bubble
	} else {
		callbacks = e.Capture
	}
	if len(callbacks) == 0 {
		ret = meta.NoListeners{}
	} else {
		ret = listenersInfo{e.mdl, e.EventModel, callbacks, capture}
		okay = true
	}
	return
}
