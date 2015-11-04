package internal

import (
	"bitbucket.org/pkg/inflect"
	"fmt"
	M "github.com/ionous/sashimi/model"
	S "github.com/ionous/sashimi/source"
	"github.com/ionous/sashimi/util/ident"
)

type PendingClasses map[ident.Id]*PendingClass // ptr for presense detection
type SingleToPlural map[string]ident.Id        // it kind of makes senses its single to plural id

type ClassFactory struct {
	allNames       NameSource
	relatives      *RelativeFactory
	pending        PendingClasses
	singleToPlural SingleToPlural
}

//
func NewClassFactory(names NameSource, rel *RelativeFactory) *ClassFactory {
	res := &ClassFactory{names, rel, make(PendingClasses), make(SingleToPlural)}
	res.addClassRef(nil, "kinds", "kind")
	res.addClassRef(nil, "data", "data")
	return res
}

//
// given the passed plural name, find the previously registered class
//
func (fac *ClassFactory) findBySingularName(singular string,
) (*PendingClass, bool) {
	id := fac.singleToPlural[singular]
	ret, okay := fac.pending[id]
	return ret, okay
}

//
// given the passed plural name, find the previously registered class
//
func (fac *ClassFactory) findByPluralName(plural string,
) (*PendingClass, bool) {
	id := M.MakeStringId(plural)
	ret, okay := fac.pending[id]
	return ret, okay
}

//
func (fac *ClassFactory) findByRelativeName(kind string, hint S.RelativeHint,
) (class *PendingClass, pluralized bool, err error) {
	switch hint & ^S.RelativeSource {
	case S.RelativeMany:
		if cls, ok := fac.findByPluralName(kind); !ok {
			err = ClassNotFound(kind)
		} else {
			class, pluralized = cls, true
		}
	case S.RelativeOne:
		if cls, ok := fac.findBySingularName(kind); !ok {
			err = ClassNotFound(kind)
		} else {
			class, pluralized = cls, false
		}
	default:
		if cls, ok := fac.findByPluralName(kind); ok {
			class, pluralized = cls, true
		} else if cls, ok := fac.findBySingularName(kind); ok {
			class, pluralized = cls, false
		} else {
			err = ClassNotFound(kind)
		}
	}
	return class, pluralized, err
}

//
func (fac *ClassFactory) makeClasses(relatives *RelativeFactory) (
	classes M.ClassMap, err error,
) {
	cr := newResults(fac.pending, relatives)
	return cr.finalizeClasses()
}

//
func (fac *ClassFactory) addClassRef(parent *PendingClass, plural, single string,
) (class *PendingClass, err error,
) {
	// FIX: sanity check singular?
	if singular, e := fac._addOptions(plural, single); e != nil {
		err = e
	} else if id, e := fac.allNames.addName(nil, plural, "class", ""); e != nil {
		err = e
	} else if class = fac.pending[id]; class != nil {
		// FIX? ratchet the class down?
		if class.parent != parent {
			err = fmt.Errorf("conflicting `%v` parent class `%v` respecified as `%v`",
				plural, class.parent, parent)
		}
	} else {
		var parentProps *PropertyBuilders
		if parent != nil {
			parentProps = &parent.props
		}
		class = &PendingClass{
			fac, parent, id, plural, singular,
			fac.allNames.NewScope(plural),
			NewProperties(parentProps),
			make(PendingRules, 0),
		}
		fac.pending[id] = class
		fac.singleToPlural[singular] = id
	}

	return class, err
}

//
// ex. name="rooms", value="room".
//
func (fac *ClassFactory) _addOptions(plural, singular string) (string, error) {
	if singular == "" {
		singular = inflect.Singularize(plural)
	}
	// reserve `room` to mean `rooms`
	// we dont return the id -- if they meant a specific singular string, we want that
	// the id is just the internals of name vs. name collision
	_, err := fac.allNames.addName(nil, singular, plural, "")
	return singular, err
}
