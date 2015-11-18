package memory

import (
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/util/ident"
)

type listenersInfo struct {
	mdl       *MemoryModel
	evt       *M.EventModel
	callbacks []M.ListenerModel
	captures  bool
}

func (l listenersInfo) NumListener() int {
	return len(l.callbacks)
}

func (l listenersInfo) ListenerNum(i int) api.Listener {
	// panics out of range
	cb := &l.callbacks[i]
	return listenerInfo{cb}
}

type listenerInfo struct {
	*M.ListenerModel
}

// GetInstance can return Empty()
func (l listenerInfo) GetInstance() ident.Id {
	return l.Instance
}

// GetClass always returns a valid class id.
func (l listenerInfo) GetClass() ident.Id {
	return l.Class
}

// GetCallback() returns a valid callback id.
func (l listenerInfo) GetCallback() ident.Id {
	return l.Callback
}

//
func (l listenerInfo) GetOptions() api.CallbackOptions {
	var opt api.CallbackOptions
	if l.UseTargetOnly() {
		opt |= api.UseTargetOnly
	}
	if l.UseAfterQueue() {
		opt |= api.UseAfterQueue
	}
	return opt
}
