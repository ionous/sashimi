package compiler

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
	S "github.com/ionous/sashimi/source"
	"github.com/ionous/sashimi/util/errutil"
	"github.com/ionous/sashimi/util/ident"
)

type MultiValueTable struct {
	cls   *PendingClass
	name  int                    // column holding a user-specified name
	remap map[int]IBuildProperty // column index => property
	count int
}

type MultiValueData struct {
	table      *MultiValueTable
	instanceId ident.Id
	values     []interface{}
}

var userNameColumn ident.Id = ident.MakeId("name")

func makeValueTable(classes *ClassFactory, class string, columns []string) (
	ret MultiValueTable, err error,
) {
	// have to delay instance data until after the instances have been created.
	if cls, ok := classes.findByPluralName(class); !ok {
		err = ClassNotFound(class)
	} else {
		ret = MultiValueTable{cls: cls, remap: make(map[int]IBuildProperty), count: len(columns)}
		missing := []string{}
		dupes := make(map[string]int)
		for idx, name := range columns {
			if d := dupes[name]; d != 0 {
				dupes[name] = d + 1
			} else {
				dupes[name] = 1
				id, idx := M.MakeStringId(name), idx+1
				// search for a column called "name"
				if id == userNameColumn {
					ret.name = idx
				} else if prop, ok := cls.props.propertyById(id); ok {
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
	}
	return ret, err
}

func (mvd *MultiValueTable) addRow(instanceFactory *InstanceFactory, code S.Code, values []interface{},
) (ret MultiValueData, err error) {
	if vcount := len(values); vcount != mvd.count {
		err = code.Errorf("mismatched columns %d values, %d columns", vcount, mvd.count)
	} else {
		var name string
		// build a valid name
		if val, ok := getByIndex(values, mvd.name-1); !ok {
			name = ident.MakeUniqueId().String()
		} else if str, ok := val.(string); ok {
			name = str
		} else {
			err = code.Errorf("name column doesnt't contain a string %T", val)
		}
		// request an instance of that name
		if inst, e := instanceFactory.addInstanceRef(mvd.cls, name, "", code); e != nil {
			err = e
		} else {
			ret = MultiValueData{mvd, inst.id, values}
		}
		// values will be merged later...
	}
	return ret, err
}

func (mvd MultiValueData) mergeInto(partials PartialInstances) (err error) {
	if inst, ok := partials.partials[mvd.instanceId]; !ok {
		err = M.InstanceNotFound(mvd.instanceId.String())
	} else {
		for col, prop := range mvd.table.remap {
			if val, ok := getByIndex(mvd.values, col-1); !ok {
				err = fmt.Errorf("unexpected column %d in %d values", col, len(mvd.values))
				break
			} else if e := inst.setProperty(prop, val); e != nil {
				err = e
				break
			}
		}
	}
	return err
}

func getByIndex(values []interface{}, index int) (ret interface{}, okay bool) {
	if index >= 0 && index < len(values) {
		ret, okay = values[index], true
	}
	return ret, okay
}
