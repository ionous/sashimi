package model

import "fmt"

// interface to access instance data.
// FIX: may need to split into hasone,hasmany subtypes
type RelativeValue struct {
	inst   *InstanceInfo
	prop   *RelativeProperty
	tables TableRelations
}

//
func (rel *RelativeValue) Property() IProperty {
	return rel.prop
}

func (rel *RelativeValue) GetRelativeProperty() *RelativeProperty {
	return rel.prop
}

//
// return the underlying value of rel variant
// FIX: need to expose a coerce to marshal the input to an acceptable value,
// and for one-to-one relations return the current value
func (rel *RelativeValue) Any() interface{} {
	return nil
}

//
// FIX: can the existance of the table be ensured before creating the relative value?
// FIX: can the tables live in the properties maybe? it doesnt really need the info here.
// or, -- since there might be a bunch of other tables -- a table reference in the property.
//
func (rel *RelativeValue) Table() (ret *TableRelation, ok bool) {
	return rel.tables.TableById(rel.prop.Relation())
}

// return a list of referenced instances
func (rel *RelativeValue) List() (ret []string) {
	// gee, so cache friendly.
	if table, ok := rel.Table(); ok {
		src := rel.inst.id.String()
		ret = table.List(src, rel.prop.IsRev())
	}
	return ret
}

//
func (rel *RelativeValue) String() string {
	list := rel.List()
	return fmt.Sprint(list)
}

//
//
// FIX: where and how to validate table.style?
// returns list of items cleared
func (rel *RelativeValue) ClearReference() (ret string, err error) {
	if rel.prop.ToMany() {
		err = fmt.Errorf("setting an object, but relation is a list")
	} else {
		if table, ok := rel.Table(); !ok {
			err = fmt.Errorf("internal error? couldn't find table for relation %+v", rel.prop.fields)
		} else {
			src := rel.inst.id.String()
			// FIX: some sort of early return.
			table.Remove(func(x, y string) (removed bool) {
				if !rel.prop.IsRev() && src == x {
					ret = y
					removed = true
				} else if rel.prop.IsRev() && src == y {
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
func (rel *RelativeValue) SetReference(other *InstanceInfo) (removed string, err error) {
	prop := rel.prop
	if !other.Class().CompatibleWith(prop.Relates()) {
		err = fmt.Errorf("%s not compatible with %+v", other, prop.fields)
	} else {
		if table, ok := rel.Table(); !ok {
			err = fmt.Errorf("internal error? couldn't find table for relation %+v", rel.prop.fields)
		} else {
			if !prop.ToMany() {
				removed, err = rel.ClearReference()
			}
			if err == nil {
				src := rel.inst.id.String()
				dst := other.id.String()
				if prop.IsRev() {
					dst, src = src, dst
				}
				table.Add(src, dst)
			}
		}
	}
	return removed, err
}
