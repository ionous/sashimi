package datastore

import (
	"bytes"
	"encoding/binary"
	"github.com/ionous/sashimi/util/ident"
)

func readNums(buf *bytes.Buffer) (ret []float64, err error) {
	var size int64
	if e := binary.Read(buf, binary.LittleEndian, &size); e != nil {
		err = e
	} else {
		ret = make([]float64, size)
		err = binary.Read(buf, binary.LittleEndian, ret)
	}
	return
}

func readStrings(buf *bytes.Buffer) (ret []string, err error) {
	var size int64
	if e := binary.Read(buf, binary.LittleEndian, &size); e != nil {
		err = e
	} else {
		ret = make([]string, size)
		for i := int64(0); i < size; i++ {
			if line, e := buf.ReadString(0); e != nil {
				err = e
				break
			} else {
				ret[i] = line[:len(line)-1]
			}
		}
	}
	return
}

// FIX? might be good to write these as multiples of the same property
// so that they can be queried.
func readIds(buf *bytes.Buffer) (ret []ident.Id, err error) {
	if strs, e := readStrings(buf); e != nil {
		err = e
	} else {
		ret = make([]ident.Id, len(strs))
		for i, s := range strs {
			ret[i] = ident.Id(s)
		}
	}
	return
}

func writeNums(buf *bytes.Buffer, nums []float64) (err error) {
	if e := binary.Write(buf, binary.LittleEndian, int64(len(nums))); e != nil {
		err = e
	} else {
		err = binary.Write(buf, binary.LittleEndian, nums)
	}
	return
}

func writeStrings(buf *bytes.Buffer, strs []string) (err error) {
	if e := binary.Write(buf, binary.LittleEndian, int64(len(strs))); e != nil {
		err = e
	} else {
		for _, s := range strs {
			buf.WriteString(s)
			buf.WriteByte(0)
		}
	}
	return
}

// FIX? might be good to write these as multiples of the same property
// so that they can be queried.
func writeIds(buf *bytes.Buffer, ids []ident.Id) error {
	strs := make([]string, len(ids))
	for i, v := range ids {
		strs[i] = string(v)
	}
	return writeStrings(buf, strs)
}
