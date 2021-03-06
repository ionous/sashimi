package call

import (
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/util/ident"
)

// MakeMarkerStorage implements a simple callback compiler and lookup
func MakeMarkerStorage() MarkerStorage {
	return MarkerStorage{
		Config{},
		make(map[ident.Id]G.Callback),
		make(map[string]int),
	}
}

type MarkerStorage struct {
	Config
	callbacks  map[ident.Id]G.Callback
	iterations map[string]int
}

// Compile provides the call.Compiler interface
func (m MarkerStorage) CompileCallback(cb G.Callback) (ret ident.Id, err error) {
	marker := m.MakeMarker(cb)
	it := marker.String()
	cnt := m.iterations[it]
	marker.Iteration, m.iterations[it] = cnt, cnt+1
	if id, e := marker.Encode(); e != nil {
		err = e
	} else {
		m.callbacks[id] = cb
		ret = id
	}
	return
}

// LookupCallback provides, via duck-typing, the runtime.Callbacks interface
func (m MarkerStorage) LookupCallback(id ident.Id) (ret G.Callback, okay bool) {
	if r, ok := m.callbacks[id]; ok {
		ret, okay = r, ok
	}
	return
}
