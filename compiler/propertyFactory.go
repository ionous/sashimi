package compiler

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
	S "github.com/ionous/sashimi/source"
)

const (
	TextType = "text"
	NumType  = "num"
	EnumType = "enum"
)

//
//
//
type PendingProperties map[M.StringId]M.IProperty

//
func (this *PendingClass) addPrimitive(src S.Code, fields S.PropertyFields,
) (prop M.IProperty, err error) {
	name, kind := fields.Name, fields.Kind
	// by using name->type this ensures that if the name existed, it is of the same type now
	// this does not exclude the same name from being used twice in the same class/property hierarchy
	// that is determined separately, after the hierarchy is known.
	if id, e := this.names.addRef(name, kind, src); e != nil {
		err = e
	} else {
		switch kind {
		case TextType:
			prop = this.props[id]
			if prop == nil {
				prop = M.NewTextProperty(id, name)
				this.props[id] = prop
			}
		case NumType:
			prop = this.props[id]
			if prop == nil {
				prop = M.NewNumProperty(id, name)
				this.props[id] = prop
			}
		}
	}
	return prop, err
}

//
//
//
func (this *PendingClass) addEnum(name string,
	choices []string,
	expects []S.PropertyExpectation,
) (prop M.IProperty, err error,
) {
	if enum, e := M.CheckedEnumeration(choices); e != nil {
		err = e
	} else {
		if id, e := this.names.addName(name, EnumType); e != nil {
			err = e
		} else {
			if this.props[id] != nil {
				err = fmt.Errorf("enumeration specified ore than once")
			} else {
				prop := M.NewEnumProperty(id, name, *enum)
				for _, expect := range expects {
					rule := PropertyRule{id, expect}
					this.rules = append(this.rules, rule)
				}
				this.props[id] = prop
			}
		}
	}
	return prop, err
}

//
//
//
func (this *PendingClass) makePropertySet() (props M.PropertySet) {
	props = make(M.PropertySet)
	for id, pending := range this.props {
		props[id] = pending
	}
	return props
}

//
//
//
func (this *PendingClass) addRelative(fields S.RelativeFields, src S.Code,
) (err error) {
	if cls, isMany, e := this.owner.findByRelativeName(fields.RelatesTo, fields.Hint); e != nil {
		err = e
	} else {
		name := fields.Property
		if id, e := this.names.addRef(name, "relation", src); e != nil {
			err = e
		} else {
			rel := PendingRelative{this, id, name, cls, fields.Relation, isMany, fields.Hint.IsReverse()}
			if old, existed := this.relatives[id]; existed {
				if old.rel != rel {
					err = fmt.Errorf("relation redefined. was %s, now %s", old, rel)
				}
			} else {
				this.relatives[id] = PendingRelativeEntry{src, rel}
			}
		}
	}
	return err
}
