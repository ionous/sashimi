package internal

import (
	M "github.com/ionous/sashimi/compiler/xmodel"
	"github.com/ionous/sashimi/util/ident"
	"strings"
)

//
// IBuildProperty lets the compiler set values for a property.
//
type IBuildProperty interface {
	SetProperty(PropertyContext) error
	BuildProperty() (M.IProperty, error)
}

//
// PropertyContext contains all of the state information necessary to create a property.
// FIX: refs and values are probably unnecessary; as is possbly class -- providing
// the compiler publishes a model which pushes out all data via tabless.
//
type PropertyContext struct {
	inst   ident.Id       // owner instance id
	tables TableRelations // source of relation data
	class  *M.ClassInfo   // finalized class; the class comes after the builder, so we dont normally have access to it.
	values PendingValues  // accumulates the object's initial values
	refs   PartialMap     // verification for the existance of other instances
	value  interface{}    // value to set to the accumulating object values
}

//
// The in-progress properties of a single class
//
type PropertyBuilders struct {
	parent *PropertyBuilders
	props  map[ident.Id]IBuildProperty
	names  map[string]ident.Id
}

func NewProperties(parent *PropertyBuilders) PropertyBuilders {
	return PropertyBuilders{parent, make(map[ident.Id]IBuildProperty), make(map[string]ident.Id)}
}

//
// Make a new property, or validate an existing one using the passed callbacks for the id'd property.
//
func (b *PropertyBuilders) make(
	id ident.Id,
	name string,
	validator func(IBuildProperty) error,
	creator func() (IBuildProperty, error),
) (
	ret IBuildProperty,
	err error,
) {
	if old, existed := b.props[id]; !existed {
		if p, e := creator(); e != nil {
			err = e
		} else {
			b.props[id] = p
			b.names[name] = id
			ret = p
		}
	} else {
		if validator != nil {
			err = validator(old)
		}
		if err == nil {
			ret = old
		}
	}
	return ret, err
}

//
// FindProperty returns the named property, searching upwards through the property hierarchy.
//
func (b *PropertyBuilders) findProperty(name string) (ret IBuildProperty, okay bool) {
	for k, v := range b.names {
		if strings.EqualFold(name, k) {
			ret, okay = b.props[v]
			break
		}
	}
	if !okay && b.parent != nil {
		ret, okay = b.parent.findProperty(name)
	}
	return
}

//
// GetProperty returns the id'd property, searching upwards through the property hierarchy.
//
func (b *PropertyBuilders) propertyById(id ident.Id) (IBuildProperty, bool) {
	prop, okay := b.props[id]
	if !okay && b.parent != nil {
		prop, okay = b.parent.propertyById(id)
	}
	return prop, okay
}
