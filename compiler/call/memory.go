package call

import (
	G "github.com/ionous/sashimi/game"
	M "github.com/ionous/sashimi/model"
)

// MakeMemoryStorage implements a simple callback compiler and lookup
func MakeMemoryStorage() MemoryStorage {
	return MemoryStorage{
		Config{},
		make(map[M.Callback]G.Callback),
		make(map[string]int),
	}
}

type MemoryStorage struct {
	Config
	callbacks  map[M.Callback]G.Callback
	iterations map[string]int
}

// Compile provides the call.Compiler interface
func (m MemoryStorage) Compile(cb G.Callback) (Marker, error) {
	marker := m.MakeMarker(cb)
	it := marker.String()
	cnt := m.iterations[it]
	marker.Iteration, m.iterations[it] = cnt, cnt+1
	m.callbacks[marker.Callback] = cb
	return marker, nil
}

// Lookup provides, via duck-typing, the runtime.Callbacks interface
func (m MemoryStorage) Lookup(marker M.Callback) G.Callback {
	return m.callbacks[marker]
}
