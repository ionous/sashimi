package compiler

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
	S "github.com/ionous/sashimi/source"
	"github.com/ionous/sashimi/util/errutil"
)

type PendingClass struct {
	classes  *ClassFactory
	parent   *PendingClass
	id       M.StringId
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
func (cls *PendingClass) addPrimitive(src S.Code, fields S.PropertyFields,
) (ret IBuildProperty, err error) {
	name, kind := fields.Name, fields.Kind
	// by using name->type cls ensures that if the name existed, it is of the same type now
	// cls does not exclude the same name from being used twice in the same class/property hierarchy
	// that is determined separately, after the hierarchy is known.
	if id, e := cls.names.addRef(name, kind, src); e != nil {
		err = e
	} else {
		ret, err = cls.props.make(id,
			func(_ IBuildProperty) error { return nil },
			func() (prop IBuildProperty, err error) {
				switch kind {
				case "text":
					prop = NewTextBuilder(id, name)
				case "num":
					prop = NewNumBuilder(id, name)
				default:
					err = fmt.Errorf("unknown property type: %s", kind)
				}
				return
			})
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
			func(_ IBuildProperty) error { return EnumMultiplySpecified(cls.id, id) },
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

//
// Add a property capable of tracking another class.
//
func (cls *PendingClass) addPointer(fields S.PropertyFields, src S.Code,
) (ret IBuildProperty,
	err error,
) {
	name, kind := fields.Name, fields.Kind
	// note: it probably doesnt make sense to allow a ratcheting down of cls
	// such a ratcheting would *increasing* restrictiveness, not permissability.
	// for example: "pointer","kind" could store "teddy bears",
	// but changed to "pointer","adult white male" could only store "teddy roosevelt"
	if id, e := cls.names.addRef(name, kind, src); e != nil {
		err = e
	} else if other, ok := cls.classes.findBySingularName(kind); !ok {
		err = ClassNotFound(kind)
	} else {
		ret, err = cls.props.make(id,
			func(old IBuildProperty) (err error) {
				if old, existed := old.(PointerBuilder); existed {
					if old.class != other.id {
						err = fmt.Errorf("pointer redefined. was %s, now %s", old.class, other.id)
					}
				}
				return
			},
			func() (IBuildProperty, error) {
				return NewPointerBuilder(id, name, other.id), nil
			})
	}
	return ret, err
}
