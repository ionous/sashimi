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

//
func (this *RelativeValue) IsMany() bool {
	return this.prop.ToMany()
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
	return this.inst.refs.tables.TableById(this.prop.relation)
}

//
func (this *RelativeValue) List() (ret []string) {
	// gee, so cache friendly.
	if table, ok := this.Table(); ok {
		src := this.inst.id.String()
		ret = table.List(src, this.prop.isRev)
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
		err = fmt.Errorf("relative values require a string %s, %#v", this.prop.id, value)
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
	prop := this.prop
	if !other.CompatibleWith(this.prop.relates) {
		err = fmt.Errorf("%s not compatible with %s", other, prop)
	} else {
		if table, ok := this.Table(); !ok {
			err = fmt.Errorf("internal error? couldn't find table for relation %v", this.prop)
		} else {
			src, dst := this.inst.id, other.inst.id
			if prop.isRev {
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
func (this *RelativeValue) ClearReference() (err error) {
	if this.IsMany() {
		err = fmt.Errorf("setting an object, but relation is a list")
	} else {
		if table, ok := this.Table(); !ok {
			err = fmt.Errorf("internal error? couldn't find table for relation %v", this.prop)
		} else {
			src := this.inst.id.String()
			table.Remove(func(x, y string) bool {
				return (!this.prop.isRev && src == x) || (this.prop.isRev && src == y)
			})
		}
	}
	return err
}

//
// FIX: table.style: where and how to validate style?
func (this *RelativeValue) SetReference(other Reference) (err error) {
	prop := this.prop
	if !other.CompatibleWith(prop.relates) {
		err = fmt.Errorf("%s not compatible with %s", other, prop)
	} else {
		if table, ok := this.Table(); !ok {
			err = fmt.Errorf("internal error? couldn't find table for relation %v", this.prop)
		} else {
			if !prop.isMany {
				err = this.ClearReference()
			}
			if err == nil {
				src := this.inst.id.String()
				dst := other.inst.id.String()
				if prop.isRev {
					dst, src = src, dst
				}
				table.Add(src, dst)
			}
		}
	}
	return err
}
