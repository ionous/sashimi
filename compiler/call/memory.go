package call

import (
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/util/ident"
)

// MakeMemoryStorage implements a simple callback compiler and lookup
func MakeMemoryStorage() MemoryStorage {
	return MemoryStorage{
		Config{},
		make(map[ident.Id]G.Callback),
		make(map[string]int),
	}
}

type MemoryStorage struct {
	Config
	callbacks  map[ident.Id]G.Callback
	iterations map[string]int
}

// Compile provides the call.Compiler interface
func (m MemoryStorage) CompileCallback(cb G.Callback) (ret ident.Id, err error) {
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

// Lookup provides, via duck-typing, the runtime.Callbacks interface
func (m MemoryStorage) Lookup(id ident.Id) G.Callback {
	return m.callbacks[id]
}
