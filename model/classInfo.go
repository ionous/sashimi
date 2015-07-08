package model

import "github.com/ionous/sashimi/util/ident"

type ClassInfo struct {
	parent      *ClassInfo
	id          ident.Id
	name        string
	singular    string
	props       PropertySet // properties only for this cls
	constraints ConstraintSet
}

func NewClassInfo(
	parent *ClassInfo,
	id ident.Id,
	plural string,
	singular string,
	props PropertySet,
	constraints ConstraintSet,
) (cls *ClassInfo) {
	return &ClassInfo{parent, id, plural, singular, props, constraints}
}

//
func (cls *ClassInfo) Id() ident.Id {
	return cls.id
}

//
func (cls *ClassInfo) Name() string {
	return cls.name
}

//
func (cls *ClassInfo) Singular() string {
	return cls.singular
}

//
func (cls *ClassInfo) String() string {
	return cls.name
}

//
func (cls *ClassInfo) Parent() *ClassInfo {
	return cls.parent
}

//
func (cls *ClassInfo) Properties() PropertySet {
	return cls.props
}

//
// Returns a new property set consisting of all properties in this cls and all parents
//
func (cls *ClassInfo) AllProperties() PropertySet {
	props := make(PropertySet)
	cls._flatten(props)
	return props
}

//
//
//
func (cls *ClassInfo) FindProperty(name string) (IProperty, bool) {
	id := MakeStringId(name)
	return cls.PropertyById(id)
}

//
//
//
func (cls *ClassInfo) PropertyById(id ident.Id) (IProperty, bool) {
	prop, okay := cls.props[id]
	if !okay && cls.parent != nil {
		prop, okay = cls.parent.PropertyById(id)
	}
	return prop, okay
}

//
//
//
func (cls *ClassInfo) Constraints() ConstraintSet {
	return cls.constraints
}

//
// FindConstraint returns the contraint of the named property.
//
func (cls *ClassInfo) FindConstraint(name string) (ret IConstrain, okay bool) {
	if prop, ok := cls.FindProperty(name); ok {
		ret, okay = cls.PropertyConstraint(prop)
	}
	return ret, okay
}

//
// Currently, only enum constraints are supported.
// NOTE: Contraints can appear in any descendent cls at or below the property it constrains.
// FIX? doesnt validate property comes from this class -- perhaps switch to id?
//
func (cls *ClassInfo) PropertyConstraint(prop IProperty) (ret IConstrain, okay bool) {
	switch prop := prop.(type) {
	case *EnumProperty:
		if c, ok := cls.ConstraintById(prop.Id()); ok {
			ret = c
		} else {
			ret = NewConstraint(prop.Enumeration)
		}
		okay = true
	}
	return ret, okay
}

func (cls *ClassInfo) ConstraintById(id ident.Id) (ret IConstrain, okay bool) {
	if c, ok := cls.constraints[id]; ok {
		ret, okay = c, ok
	} else if cls.parent != nil {
		ret, okay = cls.parent.ConstraintById(id)
	}
	return ret, okay
}

//
// CompatibleWith returns true when this class can be used in situtations which require the other class.
//
func (cls *ClassInfo) CompatibleWith(other ident.Id) bool {
	return cls.Id() == other || cls.HasParent(other)
}

//
//
//
func (cls *ClassInfo) HasParent(p ident.Id) (yes bool) {
	for c := cls.Parent(); c != nil; c = c.Parent() {
		if c.Id() == p {
			yes = true
			break
		}
	}
	return yes
}

//
//
//
func (cls *ClassInfo) PropertyByChoice(choice string) (
	prop *EnumProperty,
	index int,
	ok bool,
) {
	choiceId := MakeStringId(choice)
	prop, index = cls._propertyByChoice(choiceId)
	return prop, index, prop != nil
}

func (cls *ClassInfo) _propertyByChoice(choice ident.Id) (
	prop *EnumProperty,
	index int,
) {
	for _, p := range cls.props {
		switch enum := p.(type) {
		case *EnumProperty:
			idx, err := enum.ChoiceToIndex(choice)
			if err == nil {
				prop = enum
				index = idx
			}
		}
		if prop != nil {
			break
		}
	}
	if prop == nil && cls.parent != nil {
		prop, index = cls.parent._propertyByChoice(choice)
	}
	return prop, index
}

// NOTE: does NOT check for conflicts.
// ( trying to be a little looser than normal,
// and get to the point where the model is known to be safe at creation time. )
func (cls *ClassInfo) _flatten(props PropertySet) {
	if cls.parent != nil {
		cls.parent._flatten(props)
	}
	for k, prop := range cls.props {
		props[k] = prop
	}
}
