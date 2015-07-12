package compiler

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
	S "github.com/ionous/sashimi/source"
	"github.com/ionous/sashimi/util/errutil"
	"github.com/ionous/sashimi/util/ident"
)

type PendingClass struct {
	classes  *ClassFactory
	parent   *PendingClass
	id       ident.Id
	name     string
	singular string
	names    NameScope
	props    PropertyBuilders
	rules    PendingRules
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
	name, kind := fields.Name, fields.Kind
	// by using name->type cls ensures that if the name existed, it is of the same type now
	// cls does not exclude the same name from being used twice in the same class/property hierarchy
	// that is determined separately, after the hierarchy is known.
	if id, e := cls.names.addRef(name, kind, src); e != nil {
		err = e
	} else {
		switch kind {
		case "text":
			ret, err = cls.props.make(id, nil,
				func() (IBuildProperty, error) {
					return NewTextBuilder(id, name)
				})
		case "num":
			ret, err = cls.props.make(id, nil,
				func() (IBuildProperty, error) {
					return NewNumBuilder(id, name)
				})
		default:
			if other, ok := cls.classes.findBySingularName(kind); !ok {
				err = ClassNotFound(kind)
			} else {
				ret, err = cls.props.make(id,
					func(old IBuildProperty) (err error) {
						ptr, existed := old.(PointerBuilder)
						if !existed || ptr.class != other.id {
							err = fmt.Errorf("pointer redefined. was %s, now %s", old, other)
						}
						return
					},
					func() (IBuildProperty, error) {
						return NewPointerBuilder(id, name, other.id)
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
	expects []S.PropertyExpectation,
) (ret IBuildProperty,
	err error,
) {
	if id, e := cls.names.addName(name, "enum"); e != nil {
		err = e
	} else {
		ret, err = cls.props.make(id,
			func(_ IBuildProperty) error {
				return EnumMultiplySpecified(cls.id, id)
			},
			func() (prop IBuildProperty, err error) {
				return NewEnumBuilder(id, name, choices)
			})
		if err == nil {
			for _, expect := range expects {
				rule := PropertyRule{id, expect}
				cls.rules = append(cls.rules, rule)
			}
		}
	}
	return ret, err
}

//
// Add a relative property to cls class.
//
func (cls *PendingClass) addRelative(fields S.RelativeFields, src S.Code,
) (ret IBuildProperty,
	err error,
) {
	if other, isMany, e := cls.classes.findByRelativeName(fields.RelatesTo, fields.Hint); e != nil {
		err = e
	} else {
		name := fields.Property
		if id, e := cls.names.addRef(name, "relation", src); e != nil {
			err = e
		} else {
			relatives := cls.classes.relatives
			if relId, e := relatives.addName(fields.Relation, "relation"); e != nil {
				err = e
			} else {
				// create the relative property pointing to the generated relation data
				rel := M.RelativeFields{
					cls.id,
					id,       // property id in cls
					name,     // original property name
					other.id, // the other side of the relation
					relId,    // the id of the relation pair
					fields.Hint.IsReverse(),
					isMany}
				// .......
				ret, err = cls.props.make(id,
					func(old IBuildProperty) (err error) {
						if old, existed := old.(RelativeBuilder); existed {
							if old.fields != rel {
								err = fmt.Errorf("relation redefined. was %s, now %s", old, rel)
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
