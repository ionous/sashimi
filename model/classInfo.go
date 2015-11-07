package model

import "github.com/ionous/sashimi/util/ident"

type ClassInfo struct {
	Parent      *ClassInfo
	Id          ident.Id
	Plural      string
	Singular    string
	Properties  PropertySet
	Constraints ConstraintSet
}

//
func (cls *ClassInfo) String() string {
	return cls.Plural
}

// FIX: collapse all the class info methods into something using visit
// add a unit test to verify the traversal order and exit.
func (cls *ClassInfo) Visit(visitor func(cls *ClassInfo) bool) {
	if !visitor(cls) && cls.Parent != nil {
		cls.Parent.Visit(visitor)
	}
}

// AllProperties returns a new property set consisting of all properties in this cls and all parents.
func (cls *ClassInfo) AllProperties() PropertySet {
	props := make(PropertySet)
	cls._flatten(props)
	return props
}

// PropertyById searches through the class hierarchy for the property matching the passed name.
func (cls *ClassInfo) FindProperty(name string) (IProperty, bool) {
	id := MakeStringId(name)
	return cls.PropertyById(id)
}

// PropertyById searches through the class hierarchy for the property matching the passed id.
func (cls *ClassInfo) PropertyById(id ident.Id) (IProperty, bool) {
	prop, okay := cls.Properties[id]
	if !okay && cls.Parent != nil {
		prop, okay = cls.Parent.PropertyById(id)
	}
	return prop, okay
}

// CompatibleWith returns true when this class can be used in situtations which require the other class.
func (cls *ClassInfo) CompatibleWith(other ident.Id) bool {
	return cls.Id == other || cls.HasParent(other)
}

// HasParent searches through the class hierarchy for the parent matching the passed id.
func (cls *ClassInfo) HasParent(p ident.Id) (yes bool) {
	for c := cls.Parent; c != nil; c = c.Parent {
		if c.Id == p {
			yes = true
			break
		}
	}
	return yes
}

// PropertyByChoice searches through the class hiearchy for the first property containing the passed choice.
func (cls *ClassInfo) PropertyByChoice(choice string) (
	prop EnumProperty,
	index int,
	ok bool,
) {
	choiceId := MakeStringId(choice)
	return cls.PropertyByChoiceId(choiceId)
}

//
func (cls *ClassInfo) PropertyByChoiceId(choice ident.Id) (
	prop EnumProperty,
	index int,
	ok bool,
) {
	for _, p := range cls.Properties {
		switch enum := p.(type) {
		case EnumProperty:
			idx, err := enum.ChoiceToIndex(choice)
			if err == nil {
				prop = enum
				index = idx
				ok = true
			}
		}
		if ok {
			break
		}
	}
	if !ok && cls.Parent != nil {
		prop, index, ok = cls.Parent.PropertyByChoiceId(choice)
	}
	return prop, index, ok
}

// NOTE: does NOT check for conflicts.
// ( trying to be a little looser than normal,
// and get to the point where the model is known to be safe at creation time. )
func (cls *ClassInfo) _flatten(props PropertySet) {
	if cls.Parent != nil {
		cls.Parent._flatten(props)
	}
	for k, prop := range cls.Properties {
		props[k] = prop
	}
}
