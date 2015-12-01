package datastore

import (
	A "appengine"
	D "appengine/datastore"
	"fmt"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
)

type KeyValues struct {
	mdl meta.Model
	KeyGen
	objects LoadSavers
	ctx     A.Context
}

type LoadSavers map[ident.Id]*LoadSaver

func (ls LoadSavers) Save(tc A.Context) error {
	size := len(ls) // this can be as large as the number of Get(s)
	keys, loadSavers := make([]*D.Key, 0, size), make([]D.PropertyLoadSaver, 0, size)
	for _, loadSaver := range ls {
		if loadSaver.changed {
			keys = append(keys, loadSaver.key)
			loadSavers = append(loadSavers, loadSaver)
		}
	}
	tc.Infof("saving %d keys", len(keys))
	_, err := D.PutMulti(tc, keys, loadSavers)
	return err
}

func (kvs *KeyValues) Reset() {
	kvs.objects = make(LoadSavers)
}

func (kvs *KeyValues) Save() (err error) {
	return kvs.objects.Save(kvs.ctx)
}

// GetValue succeeds if SetValue was called on a corresponding obj.field.
func (kvs *KeyValues) GetValue(obj, field ident.Id) (ret interface{}, okay bool) {
	var err error
	inst, ok := kvs.objects[obj]
	if !ok {
		inst, err = kvs.Cache(obj)
	}
	if err != nil {
		panic(err)
	} else {
		ret, okay = inst.GetValue(field)
	}
	return
}

// SetValue always succeeds, storing the passed value to the map at obj.field.
func (kvs *KeyValues) SetValue(obj, field ident.Id, value interface{}) (err error) {
	inst, ok := kvs.objects[obj]
	if !ok {
		inst, err = kvs.Cache(obj)
	}
	if err != nil {
		panic(err)
	} else {
		err = inst.SetValue(field, value)
	}
	return
}

func (kvs *KeyValues) Cache(obj ident.Id) (ret *LoadSaver, err error) {
	if inst, ok := kvs.mdl.GetInstance(obj); !ok {
		err = fmt.Errorf("couldnt find %s", obj)
	} else {
		key := kvs.NewKey(inst)
		// verbose -- FIX: memcache first?
		//kvs.ctx.Debugf("caching %s", key)
		ls := NewLoadSaver(inst, key)
		if e := D.Get(kvs.ctx, key, ls); e != nil && e != D.ErrNoSuchEntity {
			err = e
		} else {
			kvs.objects[obj] = ls
			ret = ls
		}
	}
	return
}
