package model

import "fmt"

// interface to access instance data.
type RelativeValue struct {
	inst *InstanceInfo
	prop *RelativeProperty
}

//
func (this *RelativeValue) Property() IProperty {
	return this.prop
}

func (this *RelativeValue) GetRelativeProperty() *RelativeProperty {
	return this.prop
}

//
// return the underlying value of this variant
// FIX: always returns false for the compiler setKeyValue() change detection pattern;
// may need to expose a coerce to marshal the input to an acceptable value,
// and for one-to-one relations return the current value; then the compiler could check for changes properly.
func (this *RelativeValue) Any() (interface{}, bool) {
	return nil, false
}

//
// FIX: can the existance of the table be ensured before creating the relative value?
func (this *RelativeValue) Table() (ret *TableRelation, ok bool) {
	return this.inst.refs.tables.TableById(this.prop.Relation())
}

// return a list of referenced instances
func (this *RelativeValue) List() (ret []string) {
	// gee, so cache friendly.
	if table, ok := this.Table(); ok {
		src := this.inst.id.String()
		ret = table.List(src, this.prop.IsRev())
	}
	return ret
}

//
func (this *RelativeValue) String() string {
	list := this.List()
	return fmt.Sprint(list)
}

//
// set the value of the variant
func (this *RelativeValue) SetAny(value interface{}) (err error) {
	if other, okay := value.(string); !okay {
		err = fmt.Errorf("relative values require a string %s, %#v", this.prop.Id(), value)
	} else if other, e := this.inst.refs.FindByName(other); e != nil {
		err = e
	} else {
		err = this.AddReference(other)
	}
	return err
}

//
// FIX: where and how to validate table.style:
// maybe in the table iteself since its already searching for duplicate pairs...
func (this *RelativeValue) AddReference(other Reference) (err error) {
	if !other.CompatibleWith(this.prop.Relates()) {
		err = fmt.Errorf("%s not compatible with %v", other, this.prop.fields)
	} else {
		if table, ok := this.Table(); !ok {
			err = fmt.Errorf("internal error? couldn't find table for relation %+v", this.prop.fields)
		} else {
			src, dst := this.inst.id, other.inst.id
			if this.prop.IsRev() {
				dst, src = src, dst
			}
			table.Add(src.String(), dst.String())
			//fmt.Println("!!! added", this.prop.id, src, dst, table)
		}
	}
	return err
}

//
// FIX: where and how to validate table.style?
// returns list of items cleared
func (this *RelativeValue) ClearReference() (ret string, err error) {
	if this.prop.ToMany() {
		err = fmt.Errorf("setting an object, but relation is a list")
	} else {
		if table, ok := this.Table(); !ok {
			err = fmt.Errorf("internal error? couldn't find table for relation %+v", this.prop.fields)
		} else {
			src := this.inst.id.String()
			// FIX: some sort of early return.
			table.Remove(func(x, y string) (removed bool) {
				if !this.prop.IsRev() && src == x {
					ret = y
					removed = true
				} else if this.prop.IsRev() && src == y {
					ret = x
					removed = true
				}
				return removed
			})
		}
	}
	return ret, err
}

//
// FIX: table.style: where and how to validate style?
func (this *RelativeValue) SetReference(other Reference) (removed string, err error) {
	prop := this.prop
	if !other.CompatibleWith(prop.Relates()) {
		err = fmt.Errorf("%s not compatible with %+v", other, prop.fields)
	} else {
		if table, ok := this.Table(); !ok {
			err = fmt.Errorf("internal error? couldn't find table for relation %+v", this.prop.fields)
		} else {
			if !prop.ToMany() {
				removed, err = this.ClearReference()
			}
			if err == nil {
				src := this.inst.id.String()
				dst := other.inst.id.String()
				if prop.IsRev() {
					dst, src = src, dst
				}
				table.Add(src, dst)
			}
		}
	}
	return removed, err
}
