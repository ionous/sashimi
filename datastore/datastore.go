package datastore

import (
	A "appengine"
	"github.com/ionous/sashimi/compiler/metal"
	M "github.com/ionous/sashimi/compiler/model"
	"github.com/ionous/sashimi/meta"
)

func NewDataStore(ctx A.Context, m *M.Model) DataStore {
	// yuck!  mdl uses kvs ( for value lookup ), kvs uses mdl  (for keycreation and the load saver objects); if we shadowed the meta, we could avoid this.
	kvs := &KeyValues{}
	mdl := metal.NewMetal(m, kvs)
	kvs.mdl = mdl
	kvs.KeyGen = NewKeyGen(ctx, nil)
	kvs.ctx = ctx
	kvs.Reset()
	return DataStore{kvs, mdl}
}

// Flush writes any pending changes to the datastore.
func (ds *DataStore) Flush() error {
	return ds.kvs.Save()
}

// Drop clears the local cache of fetched data.
func (ds *DataStore) Drop() {
	ds.kvs.Reset()
}

type DataStore struct {
	kvs *KeyValues
	mdl meta.Model
}
