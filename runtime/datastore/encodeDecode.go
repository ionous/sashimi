package datastore

import (
	"bytes"
	"fmt"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
)

func Encode(ptype meta.PropertyType, val interface{}) (ret interface{}, err error) {
	switch ptype {
	case meta.NumProperty:
		if v, ok := val.(float32); !ok {
			err = fmt.Errorf("property not a float %T(%v)", val, val)
		} else {
			ret = float64(v)
		}
	case meta.TextProperty:
		if v, ok := val.(string); !ok {
			err = fmt.Errorf("property not a string %T(%v)", val, val)
		} else {
			ret = v
		}
	case meta.StateProperty, meta.ObjectProperty:
		if v, ok := val.(ident.Id); !ok {
			err = fmt.Errorf("property not an id %T(%v)", val, val)
		} else {
			ret = ident.Dash(v)
		}
	default:
		buf := new(bytes.Buffer)
		switch ptype & ^meta.ArrayProperty {
		case meta.NumProperty:
			if vs, ok := val.([]float32); !ok {
				err = fmt.Errorf("property not array of floats %T(%v)", val, val)
			} else {
				err = writeNums(buf, vs)
			}
		case meta.TextProperty:
			if vs, ok := val.([]string); !ok {
				err = fmt.Errorf("property not array of strings %T(%v)", val, val)
			} else {
				err = writeStrings(buf, vs)
			}
		case meta.ObjectProperty:
			if vs, ok := val.([]ident.Id); !ok {
				err = fmt.Errorf("property not array of ids %T(%v)", val, val)
			} else {
				err = writeIds(buf, vs)
			}
		default:
			err = fmt.Errorf("property type unknown (%v)", ptype)
		}
		if err == nil {
			ret = buf.Bytes()
		}
	}
	return
}

func Decode(ptype meta.PropertyType, val interface{}) (ret interface{}, err error) {
	switch ptype {
	case meta.NumProperty:
		if v, ok := val.(float64); !ok {
			err = fmt.Errorf("property not a float %T(%v)", val, val)
		} else {
			ret = float32(v)
		}
	case meta.TextProperty:
		if v, ok := val.(string); !ok {
			err = fmt.Errorf("property not a string %T(%v)", val, val)
		} else {
			ret = v
		}
	case meta.StateProperty, meta.ObjectProperty:
		if v, ok := val.(string); !ok {
			err = fmt.Errorf("property not an id string %T(%v)", val, val)
		} else {
			ret = ident.Id(v)
		}
	default:
		if v, ok := val.([]byte); !ok {
			err = fmt.Errorf("property not bytes %T(%v)", val, val)
		} else {
			buf := bytes.NewBuffer(v)
			switch ptype & ^meta.ArrayProperty {
			case meta.NumProperty:
				ret, err = readNums(buf)
			case meta.TextProperty:
				ret, err = readStrings(buf)
			case meta.ObjectProperty:
				ret, err = readIds(buf)
			default:
				err = fmt.Errorf("property type unknown (%v)", ptype)
			}
		}
	}
	return
}
