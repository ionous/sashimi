package datastore

import (
	"appengine"
	"appengine/datastore"
	M "github.com/ionous/sashimi/compiler/model"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/metal"
)

func NewModelStore(ctx appengine.Context, m *M.Model, parent *datastore.Key) *ModelStore {
	// yuck!  mdl uses kvs ( for value lookup ), kvs uses mdl  (for keycreation and the load saver objects); if we shadowed the meta, we could avoid this.
	kvs := &KeyValues{}
	mdl := metal.NewMetal(m, kvs)
	kvs.mdl = mdl
	kvs.KeyGen = NewKeyGen(ctx, parent)
	kvs.ctx = ctx
	kvs.Reset()
	return &ModelStore{kvs, mdl}
}

func (ds *ModelStore) Model() meta.Model {
	return ds.mdl
}

// Flush writes any pending changes to the datastore.
func (ds *ModelStore) Flush() error {
	return ds.kvs.Save()
}

// Drop clears the local cache of fetched data.
func (ds *ModelStore) Drop() {
	ds.kvs.Reset()
}

type ModelStore struct {
	kvs *KeyValues
	mdl meta.Model
}
