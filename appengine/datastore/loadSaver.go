package datastore

import (
	D "appengine/datastore"
	"fmt"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
)

// LoadSaver implements datastore.PropertyLoadSaver
type LoadSaver struct {
	meta.Instance
	key     *D.Key
	data    LoadSaverFields
	changed bool
}

type LoadSaverFields map[ident.Id]FieldValue

func NewLoadSaver(inst meta.Instance, key *D.Key) *LoadSaver {
	return &LoadSaver{inst, key, make(LoadSaverFields), false}
}

func (s *LoadSaver) GetValue(id ident.Id) (ret interface{}, okay bool) {
	if field, ok := s.data[id]; ok {
		okay, ret = true, field.Value
	}
	return
}

// note: we cant compare values at this levels
// or we get panics: ex. "comparing uncomparable type []float32"
func (s *LoadSaver) SetValue(id ident.Id, value interface{}) (err error) {
	if field, ok := s.data[id]; ok {
		field.Value = value
		s.data[id] = field // record value back; it's not a pointer.
		s.changed = true
	} else if prop, ok := s.GetProperty(id); !ok {
		err = fmt.Errorf("couldnt find property %s.%s", s, id)
	} else {
		s.data[id] = FieldValue{prop.GetType(), value}
		s.changed = true
	}
	return
}

func (s *LoadSaver) Load(ch <-chan D.Property) (err error) {
	//fmt.Println("loading", s.GetId())
	for p := range ch {
		//fmt.Println("reading", p.Name)
		propId := ident.Id(p.Name)
		if prop, ok := s.GetProperty(propId); !ok {
			err = fmt.Errorf("couldnt find property %s", propId)
			break
		} else {
			ptype := prop.GetType()
			if decoded, e := Decode(ptype, p.Value); e != nil {
				err = fmt.Errorf("decoding %s %s", prop, e)
			} else {
				s.data[propId] = FieldValue{ptype, decoded}
			}
		}
	}
	return
}

func (s *LoadSaver) Save(ch chan<- D.Property) (err error) {
	//fmt.Println("saving", s.GetId(), len(s.data))
	defer close(ch)
	// todo: write an instance level "version" property
	for id, val := range s.data {
		if encoded, e := val.Encode(); e != nil {
			err = fmt.Errorf("encoding %s %s", id, e)
			break
		} else {
			//fmt.Println("storing", id, encoded)
			ch <- D.Property{
				Name:    string(id),
				Value:   encoded,
				NoIndex: (val.Type & meta.ArrayProperty) != 0,
			}
		}
	}
	return
}
