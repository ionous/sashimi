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
	size := len(ls)
	tc.Infof("saving %d keys", size)
	keys, loadSavers := make([]*D.Key, size), make([]D.PropertyLoadSaver, size)
	idx := 0
	for _, loadSaver := range ls {
		keys[idx] = loadSaver.key
		loadSavers[idx] = loadSaver
		idx++
	}
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
	} else if val, ok := inst.fields[field]; ok {
		ret, okay = val.Value, true
	}
	return
}

// SetValue always succeeds, storing the passed value to the map at obj.field.
func (kvs *KeyValues) SetValue(obj, field ident.Id, value interface{}) (err error) {
	inst, ok := kvs.objects[obj]
	if !ok {
		inst, err = kvs.Cache(obj)
	}
	if err == nil {
		if prop, ok := inst.GetProperty(field); !ok {
			err = fmt.Errorf("coulnt find property %s.%s", obj, field)
		} else {
			inst.fields[field] = FieldValue{prop.GetType(), value}
		}
	}
	return
}

func (kvs *KeyValues) Cache(obj ident.Id) (ret *LoadSaver, err error) {
	if inst, ok := kvs.mdl.GetInstance(obj); !ok {
		err = fmt.Errorf("couldnt find %s", obj)
	} else {
		key := kvs.NewKey(inst)
		kvs.ctx.Infof("caching %s", key)
		ls := &LoadSaver{inst, key, make(LoadSaverFields)}
		if e := D.Get(kvs.ctx, key, ls); e != nil && e != D.ErrNoSuchEntity {
			err = e
		} else {
			kvs.objects[obj] = ls
			ret = ls
		}
	}
	return
}
