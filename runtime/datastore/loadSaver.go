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
	key    *D.Key
	fields LoadSaverFields
}

type LoadSaverFields map[ident.Id]FieldValue

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
				s.fields[propId] = FieldValue{ptype, decoded}
			}
		}
	}
	return
}

func (s *LoadSaver) Save(ch chan<- D.Property) (err error) {
	//fmt.Println("saving", s.GetId(), len(s.fields))
	defer close(ch)
	// todo: write an instance level "version" property
	for id, val := range s.fields {
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
