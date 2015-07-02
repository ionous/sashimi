package compiler

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
	S "github.com/ionous/sashimi/source"
	"github.com/ionous/sashimi/util/errutil"
	"github.com/satori/go.uuid"
)

type MultiValueTable struct {
	cls   *M.ClassInfo
	name  int
	remap map[int]M.IProperty // column index => property
	count int
}

type MultiValueData struct {
	table      *MultiValueTable
	instanceId M.StringId
	values     []interface{}
}

func MakeValueTable(class *M.ClassInfo, columns []string) (
	ret MultiValueTable, err error,
) {
	ret = MultiValueTable{cls: class, remap: make(map[int]M.IProperty), count: len(columns)}
	missing := []string{}
	props := class.AllProperties()
	dupes := make(map[string]int)
	for idx, name := range columns {
		if d := dupes[name]; d != 0 {
			dupes[name] = d + 1
		} else {
			dupes[name] = 1
			id, idx := M.MakeStringId(name), idx+1
			if id == "Name" {
				ret.name = idx
			} else if prop, ok := props[id]; ok {
				ret.remap[idx] = prop
			} else {
				missing = append(missing, name)
			}
		}
	}
	if len(missing) > 0 {
		err = fmt.Errorf("missing %s", missing)
	}
	for k, v := range dupes {
		if v > 1 {
			e := fmt.Errorf("duplicate columns %s", k)
			err = errutil.Append(err, e)
		}
	}
	return ret, err
}

func (this MultiValueData) mergeInto(instances M.InstanceMap) (err error) {
	if inst, ok := instances[this.instanceId]; !ok {
		err = M.InstanceNotFound(this.instanceId.String())
	} else {
		// FIXFIX NAME
		for col, prop := range this.table.remap {
			if val, ok := getByIndex(this.values, col-1); !ok {
				err = fmt.Errorf("unexpected column %d in %d values", col, len(this.values))
				break
			} else if e := setKeyValue(inst, prop.Name(), val); e != nil {
				err = e
				break
			}
		}
	}
	return err
}

func (this *MultiValueTable) AddRow(instanceFactory *InstanceFactory, code S.Code, values []interface{},
) (ret MultiValueData, err error) {
	if vcount := len(values); vcount != this.count {
		err = code.Errorf("mismatched columns %d values, %d columns", vcount, this.count)
	} else {
		var name string
		// build a valid name
		if val, ok := getByIndex(values, this.name-1); !ok {
			name = uuid.NewV4().String()
		} else if str, ok := val.(string); ok {
			name = str
		} else {
			err = code.Errorf("name column doesnt't contain a string %T", val)
		}
		// create an instance of that name
		if inst, e := instanceFactory.addInstanceRef(name, this.cls.Id(), S.Options{}, code); e != nil {
			err = e
		} else {
			ret = MultiValueData{this, inst.id, values}
		}
	}
	return ret, err
}

func getByIndex(values []interface{}, index int) (ret interface{}, okay bool) {
	if index >= 0 && index < len(values) {
		ret, okay = values[index], true
	}
	return ret, okay
}
