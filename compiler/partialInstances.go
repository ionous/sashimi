package compiler

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
	S "github.com/ionous/sashimi/source"
)

type PartialInstances struct {
	log       *ErrorLog
	instances M.InstanceMap
	tables    M.TableRelations
}

//
func (this *PartialInstances) makeData(choices []S.ChoiceStatement, kvs []S.KeyValueStatement,
) (ret M.InstanceMap, tables M.TableRelations, err error) {
	if e := this._addChoices(choices); e != nil {
		err = e
	} else if e := this._addKeyValues(kvs); e != nil {
		err = e
	} else {
		ret = this.instances
		tables = this.tables
	}
	return ret, tables, err
}

//
// add instance values for later processing
func (this *PartialInstances) _addChoices(choices []S.ChoiceStatement) (err error) {
	this.log.Println("adding instance choices")
	for _, choice := range choices {
		fields := choice.Fields()
		if inst, ok := this.instances.FindInstance(fields.Owner); !ok {
			err = this.log.AppendError(err, M.InstanceNotFound(fields.Owner))
		} else {
			if prop, index, ok := inst.Class().PropertyByChoice(fields.Choice); !ok {
				e := fmt.Errorf("no such choice: '%v'", choice)
				err = this.log.AppendError(err, e)
			} else {
				if e := _setKeyValue(inst, prop.Name(), index); e != nil {
					err = this.log.AppendError(err, e)
				}
			}
		}
	}
	return err
}

//
// add instance values for later processing
func (this *PartialInstances) _addKeyValues(kvs []S.KeyValueStatement) (err error) {
	this.log.Println("adding instance key values")
	for _, kv := range kvs {
		keyValue, src := kv.Fields(), kv.Source()
		if inst, ok := this.instances.FindInstance(keyValue.Owner); !ok {
			e := fmt.Errorf("instance not found %s @ %s", keyValue.Owner, src)
			err = this.log.AppendError(err, e)
		} else {
			if e := _setKeyValue(inst, keyValue.Key, keyValue.Value); e != nil {
				err = this.log.AppendError(err, e)
			}
		}
	}
	return err
}

//
func _setKeyValue(inst *M.InstanceInfo, name string, value interface{}) (err error) {
	if prop, ok := inst.ValueByName(name); !ok {
		err = fmt.Errorf("no such property %v", name)
	} else {
		if old, was := prop.Any(); !was {
			err = prop.SetAny(value)
		} else if old != value {
			err = ValueMismatch{prop.Property().Id(), old, value}
		}
	}
	return err
}
