package pack

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"github.com/ionous/sashimi/metal"
	"github.com/ionous/sashimi/util/ident"
	"strings"
)

type ObjectValuePack struct {
	Keys   []string `json:"k"`
	Values []string `json:"v"`
	Data   string   `json:"d"`
}

type indexer struct {
	unique  map[string]int
	strings []string
}

func newIndexer() *indexer {
	return &indexer{unique: make(map[string]int)}
}

func (i *indexer) add(s string) (ret int) {
	if idx, ok := i.unique[s]; ok {
		ret = idx
	} else {
		idx := len(i.strings)
		i.strings = append(i.strings, s)
		i.unique[s] = idx
		ret = idx
	}
	return
}

func writeVariant(w *bytes.Buffer, v int) error {
	b := make([]byte, 8)
	sz := binary.PutVarint(b, int64(v))
	_, err := w.Write(b[:sz])
	return err
}

type ValueType int

const (
	Num ValueType = iota
	Text
	Ident
	NumArray
	TextArray
	IdentArray
)

func writePrimitive(w *bytes.Buffer, values *indexer, v interface{}) bool {
	switch val := v.(type) {
	case float64:
		writeVariant(w, int(Num))
		binary.Write(w, binary.LittleEndian, val)
	case string:
		writeVariant(w, int(Text))
		writeVariant(w, values.add(val))
	case ident.Id:
		writeVariant(w, int(Ident))
		writeVariant(w, values.add(string(val)))
	default:
		return false
	}
	return true
}

func writeArray(w *bytes.Buffer, values *indexer, v interface{}) bool {
	switch array := v.(type) {
	case []float64:
		writeVariant(w, int(NumArray))
		writeVariant(w, len(array))
		for _, val := range array {
			binary.Write(w, binary.LittleEndian, val)
		}
	case []string:
		writeVariant(w, int(TextArray))
		writeVariant(w, len(array))
		for _, val := range array {
			writeVariant(w, values.add(string(val)))
		}
	case []ident.Id:
		writeVariant(w, int(IdentArray))
		writeVariant(w, len(array))
		for _, val := range array {
			writeVariant(w, values.add(string(val)))
		}
	default:
		return false
	}
	return true
}

func Pack(src metal.ObjectValueMap) ObjectValuePack {
	keys, values := newIndexer(), newIndexer()
	w := new(bytes.Buffer)

	for k, v := range src {
		parts := strings.Split(k, ".")
		writeVariant(w, keys.add(parts[0]))
		writeVariant(w, keys.add(parts[1]))
		if !writePrimitive(w, values, v) && !writeArray(w, values, v) {
			panic("unknown value") // + fmt.Sprintf("%v, %T ", v, v)
		}
	}
	stringed := base64.StdEncoding.EncodeToString(w.Bytes())
	return ObjectValuePack{
		keys.strings,
		values.strings,
		stringed,
	}
}

func readText(r *bytes.Buffer, values []string, fill func(v string)) (err error) {
	if v, e := binary.ReadVarint(r); e != nil {
		err = e
		panic(e)
	} else {
		fill(values[v])
	}
	return err
}

func Unpack(src ObjectValuePack) (metal.ObjectValueMap, error) {
	dst, err := make(metal.ObjectValueMap), error(nil)
	if buf, e := base64.StdEncoding.DecodeString(src.Data); e != nil {
		err = e
	} else {
		for r := bytes.NewBuffer(buf); err == nil && r.Len() != 0; {
			if a, e := binary.ReadVarint(r); e != nil {
				err = e
			} else if b, e := binary.ReadVarint(r); e != nil {
				err = e
			} else if kind, e := binary.ReadVarint(r); e != nil {
				err = e
			} else {
				parta, partb := src.Keys[a], src.Keys[b]
				key := strings.Join([]string{parta, partb}, ".")
				//
				switch ValueType(kind) {
				case Num:
					var v float64
					if e := binary.Read(r, binary.LittleEndian, &v); e != nil {
						err = e
						panic(e)
					} else {
						dst[key] = v
					}
				case Text:
					err = readText(r, src.Values, func(v string) {
						dst[key] = v
					})
				case Ident:
					err = readText(r, src.Values, func(v string) {
						dst[key] = ident.Id(v)
					})
				case NumArray:
					if l, e := binary.ReadVarint(r); e != nil {
						err = e
					} else {
						arr := make([]float64, l)
						for i := 0; i < int(l); i++ {
							var v float64
							if e := binary.Read(r, binary.LittleEndian, &v); e != nil {
								err = e
								break
							}
							arr[i] = v
						}
						dst[key] = arr
					}
				case TextArray:
					if l, e := binary.ReadVarint(r); e != nil {
						err = e
					} else {
						arr := make([]string, l)
						for i := 0; err == nil && i < int(l); i++ {
							err = readText(r, src.Values, func(v string) {
								arr[i] = v
							})
						}
						dst[key] = arr
					}
				case IdentArray:
					if l, e := binary.ReadVarint(r); e != nil {
						err = e
					} else {
						arr := make([]ident.Id, l)
						for i := 0; err == nil && i < int(l); i++ {
							err = readText(r, src.Values, func(v string) {
								arr[i] = ident.Id(v)
							})
						}
						dst[key] = arr
					}
				default:
					err = errors.New("unknown type")
				}
			}
		}
	}
	return dst, err
}
