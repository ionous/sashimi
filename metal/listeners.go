package metal

import (
	M "github.com/ionous/sashimi/compiler/model"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
)

type listenersInfo struct {
	callbacks []M.ListenerModel
	captures  bool
}

func (l *listenersInfo) NumListener() int {
	return len(l.callbacks)
}

func (l *listenersInfo) ListenerNum(i int) meta.Listener {
	// panics out of range
	cb := &l.callbacks[i]
	return &listenerInfo{cb}
}

type listenerInfo struct {
	*M.ListenerModel
}

// GetInstance can return Empty()
func (l *listenerInfo) GetInstance() ident.Id {
	return l.Instance
}

// GetClass always returns a valid class id.
func (l *listenerInfo) GetClass() ident.Id {
	return l.Class
}

// GetCallback() returns a valid callback id.
func (l *listenerInfo) GetCallback() meta.Callback {
	return l.Callback.ExecuteBlock
}

//
func (l *listenerInfo) GetOptions() meta.CallbackOptions {
	var opt meta.CallbackOptions
	if l.UseTargetOnly() {
		opt |= meta.UseTargetOnly
	}
	if l.UseAfterQueue() {
		opt |= meta.UseAfterQueue
	}
	return opt
}
