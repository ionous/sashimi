package internal

import (
	"fmt"
	M "github.com/ionous/sashimi/compiler/xmodel"
	S "github.com/ionous/sashimi/source"
	"github.com/ionous/sashimi/util/errutil"
	"github.com/ionous/sashimi/util/ident"
)

type PendingClass struct {
	classes      *ClassFactory
	parent       *PendingClass
	id           ident.Id
	name         string
	singular     string
	enums, names NameScope
	props        PropertyBuilders
	rules        PendingRules
}

//
func (cls *PendingClass) String() string {
	return cls.name
}

//
// Return a collection of the properties desired by cls class, indexed by id.
//
func (cls *PendingClass) makePropertySet() (props M.PropertySet, err error) {
	props = make(M.PropertySet)
	for id, pending := range cls.props.props {
		if prop, e := pending.BuildProperty(); e != nil {
			err = errutil.Append(err, e)
		} else {
			props[id] = prop
		}
	}
	return props, err
}

//
// Add a primitive ( text or number ) property to cls class.
//
func (cls *PendingClass) addProperty(src S.Code, fields S.PropertyFields,
) (ret IBuildProperty, err error) {
	name, kind, isMany := fields.Name, fields.Kind, fields.List
	// by using name->type cls ensures that if the name existed, it is of the same type now
	// cls does not exclude the same name from being used twice in the same class/property hierarchy
	// that is determined separately, after the hierarchy is known.
	if id, e := cls.names.addRef(name, kind, src); e != nil {
		err = e
	} else {
		id := ident.Join(cls.id, id)
		switch kind {
		case "text":
			ret, err = cls.props.make(id, name, nil,
				func() (IBuildProperty, error) {
					return NewTextBuilder(id, name, isMany)
				})
		case "num":
			ret, err = cls.props.make(id, name, nil,
				func() (IBuildProperty, error) {
					return NewNumBuilder(id, name, isMany)
				})
		default:
			if other, ok := cls.classes.findBySingularName(kind); !ok {
				err = ClassNotFound(kind)
			} else {
				ret, err = cls.props.make(id, name,
					func(old IBuildProperty) (err error) {
						ptr, existed := old.(PointerBuilder)
						if !existed || ptr.Class != other.id {
							err = fmt.Errorf("pointer redefined. was %s, now %s", old, other)
						}
						return
					},
					func() (IBuildProperty, error) {
						return NewPointerBuilder(id, name, other.id, isMany)
					})
			}
		}
	}
	return ret, err
}

//
// Add a primitive ( text or number ) property to cls class.
//
func (cls *PendingClass) addEnum(name string,
	choices []string,

) (ret IBuildProperty,
	err error,
) {
	id := ident.Join(cls.id, ident.MakeId(name))
	if _, e := cls.enums.addName(id.String(), cls.id.String()); e != nil {
		err = e
	} else {
		ret, err = cls.props.make(id, name,
			func(_ IBuildProperty) error {
				return EnumMultiplySpecified(cls.id, id)
			},
			func() (prop IBuildProperty, err error) {
				return NewEnumBuilder(id, name, choices)
			})
		// FIX: disabled constraints:
		// if err == nil {
		// 	for _, expect := range expects {
		// 		rule := PropertyRule{id, expect}
		// 		cls.rules = append(cls.rules, rule)
		// 	}
		// }
	}
	return ret, err
}

//
// Add a relative property to cls class.
//
func (cls *PendingClass) addRelative(fields S.RelativeProperty, src S.Code,
) (ret IBuildProperty,
	err error,
) {
	if other, isMany, e := cls.classes.findByRelativeName(fields.RelatesTo, fields.Hint); e != nil {
		err = SourceError(src, e)
	} else {
		name := fields.Property
		if id, e := cls.names.addRef(name, "relation", src); e != nil {
			err = SourceError(src, e)
		} else {
			relatives := cls.classes.relatives
			if relId, e := relatives.addName(fields.Relation, "relation"); e != nil {
				err = SourceError(src, e)
			} else {
				id := ident.Join(cls.id, id)
				// create the relative property pointing to the generated relation data
				rel := M.RelativeProperty{
					Id:       id,       // property id in cls
					Name:     name,     // original property name
					Relates:  other.id, // the other side of the relation
					Relation: relId,    // the id of the relation pair
					IsRev:    fields.Hint.IsReverse(),
					IsMany:   isMany}
				// .......
				ret, err = cls.props.make(id, name,
					func(old IBuildProperty) (err error) {
						if old, existed := old.(RelativeBuilder); existed {
							if old.fields != rel {
								e := fmt.Errorf("relation redefined. was %s, now %s", old, rel)
								err = SourceError(src, e)
							}
						}
						return
					},
					func() (IBuildProperty, error) {
						return NewRelativeBuilder(relatives, cls.id, id, fields.Relation, src, rel)
					})
			}
		}
	}
	return ret, err
}
